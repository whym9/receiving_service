package metrics

type Metrics interface {
	StartMetrics(addr string)
	RecordMetrics()
	AddMetrics(name, help, key string)
	Count(key string)
}
