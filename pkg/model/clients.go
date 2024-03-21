package model

import (
	promapi "github.com/prometheus/client_golang/api"
	"k8s.io/client-go/kubernetes"
	kneclient "knative.dev/client/pkg/eventing/v1"
	knsclient "knative.dev/client/pkg/serving/v1"
)

type Clients struct {
	KnsClient  knsclient.KnServingClient
	KneClient  kneclient.KnEventingClient
	K8sClient  kubernetes.Clientset
	PromClient promapi.Client
}
