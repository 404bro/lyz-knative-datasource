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

func Services(ctx context.Context, from time.Time, to time.Time, qm *model.QueryModel, agentURL string) backend.DataResponse {
	url := fmt.Sprintf("%s/services", agentURL)
	resp, err := http.Get(url)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	var services []model.Service
	err = json.Unmarshal(body, &services)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, err.Error())
	}
	response := backend.DataResponse{}
	var names []string
	var namespaces []string
	for _, service := range services {
		names = append(names, service.Name)
		namespaces = append(namespaces, service.Namespace)
	}
	serviceFrame := data.NewFrame("services")
	serviceFrame.Fields = append(serviceFrame.Fields,
		data.NewField("name", nil, names).SetConfig(&data.FieldConfig{DisplayName: "Name"}),
		data.NewField("namespace", nil, namespaces).SetConfig(&data.FieldConfig{DisplayName: "Namespace"}),
	)
	response.Frames = append(response.Frames, serviceFrame)
	return response
}
