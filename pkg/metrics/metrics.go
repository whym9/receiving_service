package metrics

type Metrics interface {
	StartMetrics(addr string)
	RecordMetrics()
}
