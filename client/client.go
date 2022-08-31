package main

import (
	"flag"
	"fmt"

	sender "github.com/whym9/receiving_service/pkg/sender/HTTP"
)

func main() {
	addr := *flag.String("addr", "http://localhost:8080/", "server address")
	name := *flag.String("name", "lo.pcapng", "name of the file")
	ch := make(chan []byte)

	handler := sender.HTTP_Handler{}

	go handler.StartServer(addr, &ch)

	ch <- []byte(name)

	res := <-ch

	fmt.Println(string(res))

}
