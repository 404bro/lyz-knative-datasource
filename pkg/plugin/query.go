package plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/lyz/knative/pkg/handler"
	"github.com/lyz/knative/pkg/model"
)

func loadQueryModel(source backend.DataQuery) (*model.QueryModel, error) {
	qm := model.QueryModel{}
	err := json.Unmarshal(source.JSON, &qm)
	if err != nil {
		return nil, err
	}
	return &qm, nil
}

func query(_ context.Context, query backend.DataQuery, clients model.Clients) backend.DataResponse {
	qm, err := loadQueryModel(query)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}
	switch qm.Type {
	case model.Overview:
		return handler.Overview(qm, clients)
	}
	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	return backend.DataResponse{}
}
