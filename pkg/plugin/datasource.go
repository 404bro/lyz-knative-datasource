package plugin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/lyz/knative/pkg/model"
	promapi "github.com/prometheus/client_golang/api"
	promv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	kneclient "knative.dev/client/pkg/eventing/v1"
	knsclient "knative.dev/client/pkg/serving/v1"
	kneclientv1 "knative.dev/eventing/pkg/client/clientset/versioned/typed/eventing/v1"
	knsclientv1 "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

// Make sure Datasource implements required interfaces. This is important to do
// since otherwise we will only get a not implemented error response from plugin in
// runtime. In this example datasource instance implements backend.QueryDataHandler,
// backend.CheckHealthHandler interfaces. Plugin should not implement all these
// interfaces - only those which are required for a particular task.
var (
	_ backend.QueryDataHandler      = (*Datasource)(nil)
	_ backend.CheckHealthHandler    = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

func loadPluginSettings(source backend.DataSourceInstanceSettings) (*model.PluginSettings, error) {
	settings := model.PluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, err
	}
	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)
	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *model.SecretPluginSettings {
	return &model.SecretPluginSettings{
		K8sToken: source["k8sToken"],
	}
}

// NewDatasource creates a new datasource instance.
func NewDatasource(_ context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	pluginSettings, err := loadPluginSettings(settings)
	if err != nil {
		return nil, err
	}
	k8sConfig := &rest.Config{
		Host:        pluginSettings.K8sUrl,
		BearerToken: pluginSettings.Secrets.K8sToken,
	}
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	knsClientv1, err := knsclientv1.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	knsClient := knsclient.NewKnServingClient(knsClientv1, "")
	kneClientv1, err := kneclientv1.NewForConfig(k8sConfig)
	if err != nil {
		return nil, err
	}
	kneClient := kneclient.NewKnEventingClient(kneClientv1, "")
	promClient, err := promapi.NewClient(promapi.Config{
		Address: pluginSettings.PromUrl,
	})
	if err != nil {
		return nil, err
	}
	return &Datasource{
		settings: *pluginSettings,
		clients: model.Clients{
			K8sClient:  *k8sClient,
			KnsClient:  knsClient,
			KneClient:  kneClient,
			PromClient: promClient,
		},
	}, nil
}

// Datasource is an example datasource which can respond to data queries, reports
// its health and has streaming skills.
type Datasource struct {
	settings model.PluginSettings
	clients  model.Clients
}

// Dispose here tells plugin SDK that plugin wants to clean up resources when a new instance
// created. As soon as datasource settings change detected by SDK old datasource instance will
// be disposed and a new one will be created using NewSampleDatasource factory function.
func (d *Datasource) Dispose() {
	// Clean up datasource instance resources.
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifier).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	// create response struct
	response := backend.NewQueryDataResponse()

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := query(ctx, q, d.clients)
		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}
	return response, nil
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (d *Datasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	// Check Kubernetes
	_, err := d.clients.K8sClient.CoreV1().Services("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Failed to get services: " + err.Error(),
		}, nil
	}

	// Check Prometheus
	promv1Api := promv1.NewAPI(d.clients.PromClient)
	_, _, err = promv1Api.Query(context.Background(), "up", time.Now())
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Error querying Prometheus: " + err.Error(),
		}, nil
	}

	// Check Knative Serving
	_, err = d.clients.KnsClient.ListServices(ctx)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Error querying Knative Serving: " + err.Error(),
		}, nil
	}

	// Check Knative Eventing
	_, err = d.clients.KneClient.ListBrokers(ctx)
	if err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: "Error querying Knative Eventing: " + err.Error(),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "Data source is working",
	}, nil
}
