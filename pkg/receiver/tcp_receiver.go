package receiver

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"receiving_service/internal/metrics"
)

type TCP_Handler struct {
	transferrer *chan []byte
}

func NewTCPHandler(ch *chan []byte) TCP_Handler {
	return TCP_Handler{ch}
}

func (t TCP_Handler) StartServer(addr string) {
	server, err := net.Listen("tcp", addr)
	if err != nil {

		log.Fatal(err)
		return
	}
	fmt.Println("TCP Server has started")
	go metrics.PromoHandler{}.StartMetrics(addr)
	for {
		connect, err := server.Accept()

		if err != nil {

			log.Fatal(err)
			return
		}
		go t.Receive(connect)
	}
}

func (t TCP_Handler) Receive(connect net.Conn) {
	metrics.PromoHandler{}.RecordMetrics()

	fileConntent := []byte{}

	for {
		read, err := ReceiveALL(connect, 8)

		if err != nil {

			log.Fatal(err)
			return
		}

		size := binary.BigEndian.Uint64(read)

		read, err = ReceiveALL(connect, size)

		if err != nil {

			log.Fatal(err)
			return
		}

		if size == 4 && string(read) == "STOP" {

			break
		}

		fileConntent = append(fileConntent, read...)

		fmt.Printf("File size: %v\n", size)

	}
	fmt.Println("Stopped receiving")

	*t.transferrer <- fileConntent

	statistics := <-*t.transferrer
	connect.Write([]byte(statistics))
	connect.Close()
	fmt.Println("File receiving has ended")
	fmt.Println()

}

func ReceiveALL(connect net.Conn, size uint64) ([]byte, error) {
	read := make([]byte, size)
	fmt.Println(size)
	_, err := io.ReadFull(connect, read)
	if err != nil {

		log.Fatal(err)
		return []byte{}, err
	}

	return read, nil
}
