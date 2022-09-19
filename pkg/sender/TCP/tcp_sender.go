package TCP

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/whym9/receiving_service/pkg/metrics"
)

var (
	name1 = "RabbitMQ_sent_processed_opts_total"
	help1 = "The total number of sending requsets"

	name2 = "RabbitMQ_sending_processed_errors_total"
	help2 = "The total number of sender errors"

	sent   prometheus.Counter
	errors prometheus.Counter
)

type TCP_Handler struct {
	metrics metrics.Metrics
	ch      chan []byte
}

func NewTCPHandler(m metrics.Metrics, ch chan []byte) TCP_Handler {
	return TCP_Handler{metrics: m, ch: ch}
}

func (t TCP_Handler) StartServer() {
	addr := os.Getenv("TCP_SENDER")

	sent = promauto.NewCounter(prometheus.CounterOpts{
		Name: name1,
		Help: help1,
	})
	errors = promauto.NewCounter(prometheus.CounterOpts{
		Name: name2,
		Help: help2,
	})

	connect, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	defer connect.Close()
	func() {
		var file []byte
		for {
			file = <-t.ch
			fmt.Println(len(file))
			name, err := t.Upload(file, connect)
			if err != nil {
				name = []byte("could not make statistics.")
				t.ch <- name
				break
			}

			t.ch <- name
		}

	}()

}

func (t TCP_Handler) Upload(file []byte, connect net.Conn) ([]byte, error) {
	sent.Inc()
	be := 0
	en := 1024

	for {

		if en > len(file) {
			bin := make([]byte, 8)
			binary.BigEndian.PutUint64(bin, uint64(len(file)-be))

			if _, err := connect.Write(bin); err != nil {
				errors.Inc()
				return []byte{}, err
			}

			if _, err := connect.Write(file[be:]); err != nil {
				errors.Inc()
				return []byte{}, err
			}
			bin = make([]byte, 8)

			binary.BigEndian.PutUint64(bin, uint64(4))

			if _, err := connect.Write(bin); err != nil {
				errors.Inc()
				return []byte{}, err
			}

			if _, err := connect.Write([]byte("STOP")); err != nil {
				errors.Inc()
				return []byte{}, err
			}

			break
		}
		bin := make([]byte, 8)
		binary.BigEndian.PutUint64(bin, uint64(1024))

		if _, err := connect.Write(bin); err != nil {
			errors.Inc()
			return []byte{}, err
		}

		if _, err := connect.Write(file[be:en]); err != nil {
			errors.Inc()
			return []byte{}, err
		}

		be = en
		en += 1024
	}

	read := make([]byte, 1024)

	_, err := connect.Read(read)

	if err != nil {
		errors.Inc()
		return []byte{}, err
	}

	connect.Close()
	return read, nil
}
