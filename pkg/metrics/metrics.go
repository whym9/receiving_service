package metrics

type Metrics interface {
	StartMetrics()
	RecordMetrics()
	AddMetrics(name, help, key string)
	Count(key string)
}
