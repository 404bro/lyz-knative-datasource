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

func query(ctx context.Context, query backend.DataQuery, agentURL string) backend.DataResponse {
	from := query.TimeRange.From
	to := query.TimeRange.To
	qm, err := loadQueryModel(query)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err.Error()))
	}
	switch qm.Type {
	case model.QueryServiceMap:
		return handler.Overview(ctx, from, to, qm, agentURL)
	case model.QueryServices:
		return handler.Services(ctx, from, to, qm, agentURL)
	case model.QueryAnalyticsMap:
		return handler.AnalyticsMap(ctx, from, to, qm, agentURL)
	case model.QueryAnalyticsServices:
		return handler.AnalyticsServices(ctx, from, to, qm, agentURL)
	case model.QueryAnalyticsTraces:
		return handler.AnalyticsTraces(ctx, from, to, qm, agentURL)
	}
	// create data frame response.
	// For an overview on data frames and how grafana handles them:
	// https://grafana.com/developers/plugin-tools/introduction/data-frames
	return backend.DataResponse{}
}
