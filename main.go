package main

import (
	"flag"
	"receiving_service/internal/receiver"
	"receiving_service/internal/sender"
)

var rabbit_addr string = "amqp://guest:guest@localhost:5672/"

func main() {
	addr := *flag.String("addr", "localhost:8080", "TCP server adress")
	saddr := *flag.String("saddr", ":5005", "GRPC address")
	ch := make(chan []byte)

	handler := receiver.NewHTTPHandler(&ch)

	go sender.Client{}.StartServer(saddr, &ch)

	handler.StartServer(addr)

}
