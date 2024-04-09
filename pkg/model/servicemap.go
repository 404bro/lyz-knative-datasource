package model

type ServiceMap struct {
	Nodes []ServiceMapNode `json:"nodes" bson:"nodes"`
	Edges []ServiceMapEdge `json:"edges" bson:"edges"`
}

type ServiceMapNode struct {
	Name      string  `json:"name" bson:"name"`
	Namespace string  `json:"namespace" bson:"namespace"`
	RPS       float64 `json:"rps" bson:"rps"`
	Latency   float64 `json:"latency" bson:"latency"`
	Success   float64 `json:"success" bson:"success"`
}

type ServiceMapEdge struct {
	SrcName      string `json:"srcName" bson:"srcName"`
	SrcNamespace string `json:"srcNamespace" bson:"srcNamespace"`
	DstName      string `json:"dstName" bson:"dstName"`
	DstNamespace string `json:"dstNamespace" bson:"dstNamespace"`
}
