package metrics

type Metrics interface {
	Metrics(addr string)
	RecordMetrics()
}
