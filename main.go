package main

import (
	"github.com/whym9/receiving_service/internal/worker"
	receiver "github.com/whym9/receiving_service/pkg/receiver/HTTP"

	metrics "github.com/whym9/receiving_service/pkg/metrics/prometheus"
	sender "github.com/whym9/receiving_service/pkg/sender/GRPC"
)

var rabbit_addr string = "amqp://guest:guest@localhost:5672/"

func main() {

	ch := make(chan []byte)

	Promo_Handler := metrics.NewPromoHandler()

	HTTP_handler := receiver.NewHTTPHandler(Promo_Handler, ch)

	GRPC_Handler := sender.NewGRPCHandler(Promo_Handler, ch)

	w := worker.NewWorker(GRPC_Handler, HTTP_handler, Promo_Handler)

	w.Work()

}
