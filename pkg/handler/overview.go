package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/lyz/knative/pkg/model"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	pmodel "github.com/prometheus/common/model"
)

func Overview(ctx context.Context, from time.Time, to time.Time, qm *model.QueryModel, clients model.Clients) backend.DataResponse {
	response := backend.DataResponse{}
	dur := to.Sub(from)
	var frameServiceUID []string
	var frameServiceName []string
	var frameServiceNamespace []string
	var frameServiceRPS []float64
	var frameServiceLatencies []float64
	var frameServiceSuccess []float64
	var frameServiceFail []float64

	nns2uid := make(map[model.NNS]string)
	hasNNS := make(map[model.NNS]bool)
	promv1Api := promv1.NewAPI(clients.PromClient)
	serviceList, err := clients.KnsClient.ListServices(ctx)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadGateway, err.Error())
	}
	for _, service := range serviceList.Items {
		nns2uid[model.NNS{Name: service.Name, Namespace: service.Namespace}] = string(service.UID)
		hasNNS[model.NNS{Name: service.Name, Namespace: service.Namespace}] = true
		frameServiceUID = append(frameServiceUID, string(service.UID))
		frameServiceName = append(frameServiceName, service.Name)
		frameServiceNamespace = append(frameServiceNamespace, service.Namespace)

		promQueryStr := fmt.Sprintf("scalar(sum(rate(revision_app_request_count{configuration_name=\"%s\"}[%ds])))",
			service.Name, int(dur.Seconds()))
		result, _, err := promv1Api.Query(ctx, promQueryStr, to)
		if err != nil {
			return backend.ErrDataResponse(backend.StatusInternal, err.Error())
		}
		value := float64(result.(*pmodel.Scalar).Value)
		if math.IsNaN(value) {
			value = 0
		}
		frameServiceRPS = append(frameServiceRPS, value)

		promQueryStr = fmt.Sprintf("scalar((sum(delta(revision_app_request_latencies_sum{configuration_name=\"%s\"} [%ds]))) "+
			"/ (sum(delta(revision_app_request_latencies_count{configuration_name=\"%s\"} [%ds]))))",
			service.Name, int(dur.Seconds()), service.Name, int(dur.Seconds()))
		result, _, err = promv1Api.Query(ctx, promQueryStr, to)
		if err != nil {
			return backend.ErrDataResponse(backend.StatusInternal, err.Error())
		}
		value = float64(result.(*pmodel.Scalar).Value)
		if math.IsNaN(value) {
			value = 0
		}
		frameServiceLatencies = append(frameServiceLatencies, value)

		promQueryStr = fmt.Sprintf("scalar((sum(delta(revision_app_request_count{configuration_name=\"%s\", "+
			"response_code_class=\"2xx\"} [%ds]))))", service.Name, int(dur.Seconds()))
		result, _, err = promv1Api.Query(ctx, promQueryStr, to)
		if err != nil {
			return backend.ErrDataResponse(backend.StatusInternal, err.Error())
		}
		successCount := float64(result.(*pmodel.Scalar).Value)
		if math.IsNaN(successCount) {
			successCount = 0
		}

		promQueryStr = fmt.Sprintf("scalar((sum(delta(revision_app_request_count{configuration_name=\"%s\", "+
			"response_code_class=\"5xx\"} [%ds]))))", service.Name, int(dur.Seconds()))
		result, _, err = promv1Api.Query(ctx, promQueryStr, to)
		if err != nil {
			return backend.ErrDataResponse(backend.StatusInternal, err.Error())
		}
		failCount := float64(result.(*pmodel.Scalar).Value)
		if math.IsNaN(failCount) {
			failCount = 0
		}
		if successCount == 0 && failCount == 0 {
			successCount = 1
		}
		frameServiceSuccess = append(frameServiceSuccess, successCount/(successCount+failCount))
		frameServiceFail = append(frameServiceFail, failCount/(successCount+failCount))
	}

	nodeFrame := data.NewFrame("nodes")
	nodeFrame.Fields = append(nodeFrame.Fields,
		data.NewField("id", nil, frameServiceUID),
		data.NewField("title", nil, frameServiceName).SetConfig(&data.FieldConfig{DisplayName: "Name"}),
		data.NewField("subtitle", nil, frameServiceNamespace).SetConfig(&data.FieldConfig{DisplayName: "Namespace"}),
		data.NewField("mainstat", nil, frameServiceRPS).SetConfig(&data.FieldConfig{DisplayName: "RPS"}),
		data.NewField("secondarystat", nil, frameServiceLatencies).SetConfig(&data.FieldConfig{DisplayName: "Latency"}),
		data.NewField("arc__success", nil, frameServiceSuccess).
			SetConfig(&data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "green"}, DisplayName: "Success"}),
		data.NewField("arc__fail", nil, frameServiceFail).
			SetConfig(&data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "red"}, DisplayName: "Fail"}),
	)

	var frameEdgeUID []string
	var frameSource []string
	var frameTarget []string
	edges, err := getEdges(from, to, clients.JaegerUrl, &hasNNS)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	for edge, ok := range edges {
		if ok {
			srcUID, ok1 := nns2uid[model.NNS{Name: edge.SrcName, Namespace: edge.SrcNamespace}]
			dstUID, ok2 := nns2uid[model.NNS{Name: edge.DstName, Namespace: edge.DstNamespace}]
			if ok1 && ok2 {
				frameEdgeUID = append(frameEdgeUID, uuid.New().String())
				frameSource = append(frameSource, srcUID)
				frameTarget = append(frameTarget, dstUID)
			} else {
				backend.Logger.Error("cannot find node from pair{name, namespace}")
			}
		}
	}

	edgeFrame := data.NewFrame("edges")
	edgeFrame.Fields = append(edgeFrame.Fields,
		data.NewField("id", nil, frameEdgeUID),
		data.NewField("source", nil, frameSource),
		data.NewField("target", nil, frameTarget),
	)

	response.Frames = append(response.Frames, nodeFrame)
	response.Frames = append(response.Frames, edgeFrame)
	return response
}

func getEdges(from time.Time, to time.Time, url string, hasNNS *map[model.NNS]bool) (map[model.Edge]bool, error) {
	spans := make(map[string]*model.Span)
	traces := make(map[string]*model.Trace)
	roots := make(map[string]bool)
	children := make(map[string][]string)
	edges := make(map[model.Edge]bool)
	serviceUrl := fmt.Sprintf("%s/api/services", url)
	serviceResp, err := http.Get(serviceUrl)
	if err != nil {
		return nil, err
	}
	defer serviceResp.Body.Close()
	serviceBody, err := io.ReadAll(serviceResp.Body)
	if err != nil {
		return nil, err
	}
	services := model.Services{}
	json.Unmarshal(serviceBody, &services)
	// process JSON data in each service
	for _, service := range services.Data {
		url := fmt.Sprintf("%s/api/traces?end=%d&limit=0&service=%s&start=%d", url, to.UnixMicro(), service, from.UnixMicro())
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var data model.Traces
		json.Unmarshal(body, &data)
		for _, trace := range data.Data {
			if _, ok := traces[trace.TraceID]; ok {
				continue
			}
			trace_copy := trace
			traces[trace.TraceID] = &trace_copy
			for _, span := range trace.Spans {
				if _, ok := spans[span.SpanID]; ok {
					continue
				}
				for _, tag := range span.Tags {
					if tag.Key == "http.url" {
						regexPattern := `^https?://([a-zA-Z0-9-]+)\.([a-zA-Z0-9-]+)\.svc\.cluster\.local$`
						regex := regexp.MustCompile(regexPattern)
						matches := regex.FindStringSubmatch(tag.Value)
						if len(matches) > 2 {
							span.ServiceName = matches[1]
							span.ServiceNamespace = matches[2]
							break
						}
					}
					if tag.Key == "http.host" {
						parts := strings.Split(tag.Value, ".")
						span.IsParallel = strings.HasSuffix(parts[0], "-kn-parallel-kn-channel") || strings.HasSuffix(parts[0], "-kne-trigger-kn-channel")
					}
				}
				span_copy := span
				spans[span.SpanID] = &span_copy
			}
		}
	}
	// set parent, children and roots
	for _, span := range spans {
		childOf := ""
		for _, ref := range span.References {
			if ref.RefType == "CHILD_OF" {
				childOf = ref.SpanID
				break
			}
		}
		if childOf == "" {
			roots[span.SpanID] = true
		} else {
			if _, ok := spans[childOf]; ok {
				children[childOf] = append(children[childOf], span.SpanID)
				span.ParentSpanID = childOf
			} else {
				roots[span.SpanID] = true
			}
		}
	}
	// sort children by start time
	for key, childrenList := range children {
		sort.Slice(childrenList, func(i, j int) bool {
			return spans[childrenList[i]].StartTime < spans[childrenList[j]].StartTime
		})
		children[key] = childrenList
	}
	// define functions
	isVisible := func(spanId string) bool {
		name := spans[spanId].ServiceName
		namespace := spans[spanId].ServiceNamespace
		return (*hasNNS)[model.NNS{Name: name, Namespace: namespace}]
	}
	isParallel := func(spanId string) bool {
		return spans[spanId].IsParallel
	}
	var dfs func(string, string) string
	dfs = func(spanName string, src string) string {
		if isVisible(spanName) {
			if src != "" {
				edges[model.Edge{SrcName: spans[src].ServiceName, SrcNamespace: spans[src].ServiceNamespace,
					DstName: spans[spanName].ServiceName, DstNamespace: spans[spanName].ServiceNamespace}] = true
			}
			src = spanName
		}
		for _, child := range children[spanName] {
			r := dfs(child, src)
			if !isParallel(spanName) {
				src = r
			}
		}
		return src
	}
	// execute DFS
	for root, isRoot := range roots {
		if isRoot {
			dfs(root, "")
		}
	}
	return edges, nil
}
