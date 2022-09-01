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

func (w worker) Work(addr1, addr2, addr3 string) {
	go w.metrics.StartMetrics(addr3)
	go w.sender.StartServer(addr2)

	w.receiver.StartServer(addr1)
}
