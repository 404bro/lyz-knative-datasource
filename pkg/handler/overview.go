package handler

import (
	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/lyz/knative/pkg/model"
)

func Overview(qm *model.QueryModel, client model.Clients) backend.DataResponse {
	response := backend.DataResponse{}
	frame := data.NewFrame("nodes")

	// add fields.
	frame.Fields = append(frame.Fields,
		data.NewField("id", nil, []string{"node-0", "node-1", "node-2"}),
		data.NewField("title", nil, []string{"Node 0", "Node 1", "Node 2"}).SetConfig(&data.FieldConfig{DisplayName: "Node"}),
		data.NewField("subtitle", nil, []string{"sub-1", "sub-2", "sub-3"}),
		data.NewField("mainstat", nil, []float64{1.23, 4.56, 7.89}),
		data.NewField("secondarystat", nil, []int64{11, 22, 33}),
		data.NewField("arc__red", nil, []float64{0.1, 0.2, 0.3}).SetConfig(&data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "red"}}),
		data.NewField("arc__green", nil, []float64{0.9, 0.8, 0.7}).SetConfig(&data.FieldConfig{Color: map[string]interface{}{"mode": "fixed", "fixedColor": "green"}}),
		data.NewField("detail__hello", nil, []string{"detail 1", "detail 2", "detail 3"}),
	)

	// add the frames to the response.
	response.Frames = append(response.Frames, frame)
	return response
}
