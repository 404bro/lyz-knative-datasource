package model

const (
	QueryServiceMap        = "Service Map"
	QueryServices          = "Services"
	QueryAnalyticsMap      = "Analytics Map"
	QueryAnalyticsTraces   = "Analytics Traces"
	QueryAnalyticsServices = "Analytics Services"
)

type QueryModel struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}
