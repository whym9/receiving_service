package main

import (
	"flag"

	"github.com/whym9/receiving_service/internal/worker"
	receiver "github.com/whym9/receiving_service/pkg/receiver/HTTP"

	metrics "github.com/whym9/receiving_service/pkg/metrics/prometheus"
	sender "github.com/whym9/receiving_service/pkg/sender/GRPC"
)

var rabbit_addr string = "amqp://guest:guest@localhost:5672/"

func main() {
	addr1 := *flag.String("addr1", "localhost:8080", "TCP server adress")
	addr2 := *flag.String("addr2", ":6006", "GRPC address")
	addr3 := *flag.String("addr3", ":8008", "metrics address")
	ch := make(chan []byte)

	HTTP_handler := receiver.NewHTTPHandler(&ch)

	GRPC_Handler := sender.NewGRPCHandler(&ch)

	Promo_Handler := metrics.NewPromoHandler()

	w := worker.NewWorker(&GRPC_Handler, HTTP_handler, Promo_Handler)

	w.Work(addr1, addr2, addr3)

}
