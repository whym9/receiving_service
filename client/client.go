package main

import (
	"flag"
	"fmt"

	metrics "github.com/whym9/receiving_service/pkg/metrics/prometheus"
	sender "github.com/whym9/receiving_service/pkg/sender/HTTP"
)

func main() {

	name := *flag.String("name", "lo.pcap", "name of the file")
	ch := make(chan []byte)
	metrics := metrics.NewPromoHandler()
	handler := sender.NewHTTPHandler(metrics, ch)

	go handler.StartServer()

	ch <- []byte(name)

	res := <-ch

	fmt.Println(string(res))

}
