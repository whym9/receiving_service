package RabbitMQ

import (
	"fmt"
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/streadway/amqp"
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

type Rabbit_Handler struct {
	metrics metrics.Metrics
	tr      chan []byte
}

func NewRabbitMQHandler(m metrics.Metrics, tr chan []byte) Rabbit_Handler {
	return Rabbit_Handler{metrics: m, tr: tr}
}

func (r Rabbit_Handler) StartServer() {
	addr := os.Getenv("RABBITMQ_SENDER")

	sent = promauto.NewCounter(prometheus.CounterOpts{
		Name: name1,
		Help: help1,
	})
	errors = promauto.NewCounter(prometheus.CounterOpts{
		Name: name2,
		Help: help2,
	})
	fmt.Println("RabbitMq!")

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ Instance")

	file := <-r.tr

	mes, err := r.Upload(file, "client", *conn)

	if err != nil {
		log.Fatal(err)
	}

	r.tr <- mes

}

func (r Rabbit_Handler) Upload(file []byte, name string, conn amqp.Connection) ([]byte, error) {
	sent.Inc()
	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	err = Publisher(ch, "Server", []byte(name))

	if err != nil {
		errors.Inc()
		return []byte{}, err
	}

	be := 0
	en := 1024
	for {

		if len(file) < en {
			err = Publisher(ch, name, file[be:])
			if err != nil {
				errors.Inc()
				return []byte{}, err
			}

			err = Publisher(ch, name, []byte("Stop"))
			if err != nil {
				errors.Inc()
				return []byte{}, err
			}

			break
		}
		err = Publisher(ch, name, file[be:en])

		if err != nil {
			errors.Inc()
			return []byte{}, err
		}

		be = en
		en += 1024
	}

	fmt.Println("Successfully Published Message to Queue")
	ch, err = conn.Channel()

	if err != nil {
		errors.Inc()
		return []byte{}, err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		name+"*",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		errors.Inc()
		return []byte{}, err
	}
	var res []byte
	for d := range msgs {
		res = d.Body
		fmt.Println("Successfully received messages")

		break
	}

	return res, nil
}

func Publisher(ch *amqp.Channel, name string, mes []byte) error {
	err := ch.Publish(
		"",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/octet-stream",
			Body:        mes,
		},
	)

	return err
}
