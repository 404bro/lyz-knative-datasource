package model

type PluginSettings struct {
	K8sUrl  string                `json:"k8sUrl"`
	PromUrl string                `json:"promUrl"`
	Secrets *SecretPluginSettings `json:"-"`
}

type SecretPluginSettings struct {
	K8sToken string `json:"k8sToken"`
}
