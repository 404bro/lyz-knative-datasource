package models

import (
	"encoding/json"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

type PluginSettings struct {
	K8sUrl  string                `json:"k8sUrl"`
	PromUrl string                `json:"promUrl"`
	Secrets *SecretPluginSettings `json:"-"`
}

type SecretPluginSettings struct {
	K8sToken string `json:"k8sToken"`
}

func LoadPluginSettings(source backend.DataSourceInstanceSettings) (*PluginSettings, error) {
	settings := PluginSettings{}
	err := json.Unmarshal(source.JSONData, &settings)
	if err != nil {
		return nil, err
	}
	settings.Secrets = loadSecretPluginSettings(source.DecryptedSecureJSONData)
	return &settings, nil
}

func loadSecretPluginSettings(source map[string]string) *SecretPluginSettings {
	return &SecretPluginSettings{
		K8sToken: source["k8sToken"],
	}
}
