package model

type Traces struct {
	Data []Trace `json:"data"`
}

type Trace struct {
	TraceID string `json:"traceID"`
	Spans   []Span `json:"spans"`
}

type Span struct {
	TraceID       string      `json:"traceID"`
	SpanID        string      `json:"spanID"`
	OperationName string      `json:"operationName"`
	References    []Reference `json:"references"`
	StartTime     int64       `json:"startTime"`
	Duration      int64       `json:"duration"`
	Tags          []Tag       `json:"tags"`

	SrcSpanID        string `json:"srcSpanID,omitempty"`
	ParentSpanID     string `json:"parentSpanID,omitempty"`
	ServiceName      string `json:"serviceName,omitempty"`
	ServiceNamespace string `json:"serviceNamespace,omitempty"`
	IsParallel       bool   `json:"isParallel,omitempty"`
}

type Reference struct {
	RefType string `json:"refType"`
	TraceID string `json:"traceID"`
	SpanID  string `json:"spanID"`
}

type Tag struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Edge struct {
	SrcName      string
	SrcNamespace string
	DstName      string
	DstNamespace string
}

type NNS struct {
	Name      string
	Namespace string
}

type Services struct {
	Data []string `json:"data"`
}
