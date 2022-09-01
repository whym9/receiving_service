package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type promoHandler struct {
}

func NewPromoHandler() promoHandler {
	return promoHandler{}
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func (p promoHandler) RecordMetrics() {

	opsProcessed.Inc()

}

func (p promoHandler) StartMetrics(addr string) {
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(addr, nil)
}
