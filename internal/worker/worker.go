package worker

import (
	"github.com/whym9/receiving_service/pkg/metrics"
	"github.com/whym9/receiving_service/pkg/receiver"
	"github.com/whym9/receiving_service/pkg/sender"
)

type worker struct {
	sender   sender.Sender
	receiver receiver.Receiver
	metrics  metrics.Metrics
}

func NewWorker(s sender.Sender, r receiver.Receiver, m metrics.Metrics) worker {
	return worker{
		sender:   s,
		receiver: r,
		metrics:  m,
	}
}

func (w worker) Work() {
	go w.metrics.StartMetrics("443")
	go w.sender.StartServer(":6006")

	w.receiver.StartServer("localhost:8080")
}
