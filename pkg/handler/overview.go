package handler

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/lyz/knative/pkg/model"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	pmodel "github.com/prometheus/common/model"
)

func Overview(ctx context.Context, from time.Time, to time.Time, qm *model.QueryModel, clients model.Clients) backend.DataResponse {
	response := backend.DataResponse{}
	dur := to.Sub(from)
	frame := data.NewFrame("nodes")
	var frameServiceUID []string
	var frameServiceName []string
	var frameServiceNamespace []string
	var frameServiceRPS []float64
	var frameServiceLatencies []float64
	var frameServiceSuccess []float64
	var frameServiceFail []float64

	promv1Api := promv1.NewAPI(clients.PromClient)
	serviceList, err := clients.KnsClient.ListServices(ctx)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadGateway, err.Error())
	}
	for _, service := range serviceList.Items {
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

	frame.Fields = append(frame.Fields,
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

	response.Frames = append(response.Frames, frame)
	return response
}
