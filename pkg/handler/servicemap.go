package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/lyz/knative/pkg/model"
)

func Overview(ctx context.Context, from time.Time, to time.Time, qm *model.QueryModel, agentURL string) backend.DataResponse {
	url := fmt.Sprintf("%s/map?from=%d&to=%d", agentURL, from.UnixMicro(), to.UnixMicro())
	resp, err := http.Get(url)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	var serviceMap model.ServiceMap
	err = json.Unmarshal(body, &serviceMap)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}

	response := backend.DataResponse{}
	var frameServiceUID []string
	var frameServiceName []string
	var frameServiceNamespace []string
	var frameServiceRPS []float64
	var frameServiceLatency []float64
	var frameServiceSuccess []float64
	var frameServiceFail []float64

	for _, node := range serviceMap.Nodes {
		frameServiceUID = append(frameServiceUID, node.Name+"_"+node.Namespace)
		frameServiceName = append(frameServiceName, node.Name)
		frameServiceNamespace = append(frameServiceNamespace, node.Namespace)
		frameServiceRPS = append(frameServiceRPS, node.RPS)
		frameServiceLatency = append(frameServiceLatency, node.Latency)
		frameServiceSuccess = append(frameServiceSuccess, node.Success)
		frameServiceFail = append(frameServiceFail, 1-node.Success)
	}

	nodeFrame := data.NewFrame("nodes")
	nodeFrame.Fields = append(nodeFrame.Fields,
		data.NewField("id", nil, frameServiceUID),
		data.NewField("title", nil, frameServiceName).SetConfig(&data.FieldConfig{DisplayName: "Name"}),
		data.NewField("subtitle", nil, frameServiceNamespace).SetConfig(&data.FieldConfig{DisplayName: "Namespace"}),
		data.NewField("mainstat", nil, frameServiceRPS).SetConfig(&data.FieldConfig{DisplayName: "RPS"}),
		data.NewField("secondarystat", nil, frameServiceLatency).SetConfig(&data.FieldConfig{DisplayName: "Latency"}),
		data.NewField("arc__success", nil, frameServiceSuccess).
			SetConfig(&data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "green"}, DisplayName: "Success"}),
		data.NewField("arc__fail", nil, frameServiceFail).
			SetConfig(&data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "red"}, DisplayName: "Fail"}),
	)

	var frameEdgeUID []string
	var frameSource []string
	var frameTarget []string

	for _, edge := range serviceMap.Edges {
		frameEdgeUID = append(frameEdgeUID, edge.SrcName+"_"+edge.SrcNamespace+"_"+edge.DstName+"_"+edge.DstNamespace)
		frameSource = append(frameSource, edge.SrcName+"_"+edge.SrcNamespace)
		frameTarget = append(frameTarget, edge.DstName+"_"+edge.DstNamespace)
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
