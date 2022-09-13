package RabbitMQ

import (
	"fmt"
	"log"

	"github.com/whym9/receiving_service/pkg/metrics"

	"github.com/streadway/amqp"
)

var (
	name = "RabbitMQ_receiver_processed_errors_total"
	help = "The total number of receiver errors"
	key  = "errors"
)

type Rabbit_Handler struct {
	metrics     metrics.Metrics
	transferrer chan []byte
}

func NewRabbitHandler(m metrics.Metrics, ch chan []byte) Rabbit_Handler {
	return Rabbit_Handler{metrics: m, transferrer: ch}
}

func (r Rabbit_Handler) StartServer(addr string) {
	r.metrics.AddMetrics(name, help, key)
	conn, err := amqp.Dial(addr)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer ch.Close()
	err = Declerer(*ch, "Server")

	if err != nil {
		log.Fatal(err)
	}

	msgs, err := Consumer(ch, "Server")
	if err != nil {
		r.metrics.Count(key)
		log.Fatal(err)
	}

	go func() {
		for d := range msgs {

			r.Receive(ch, string(d.Body))
		}
	}()

	fmt.Println("RabbitMQ server has started")

}

func (r Rabbit_Handler) Receive(ch *amqp.Channel, name string) {
	r.metrics.RecordMetrics()
	err := Declerer(*ch, name)

	if err != nil {
		r.metrics.Count(key)
		log.Fatal(err)
	}

	msgs, err := Consumer(ch, name)
	if err != nil {
		r.metrics.Count(key)
		log.Fatal(err)
	}
	rec := []byte{}

	for d := range msgs {
		if string(d.Body) == "Stop" {
			r.transferrer <- rec

			mes := string(<-r.transferrer)
			err = Publisher(ch, mes, name+"*")
			if err != nil {
				r.metrics.Count(key)
				log.Fatal(err)
			}
			fmt.Println("Successfully stopped receiving file!")
			break
		}
		rec = append(rec, d.Body...)
	}

}

func Publisher(ch *amqp.Channel, mes, name string) error {
	_, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}
	err = ch.Publish(
		"",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(mes),
		},
	)
	return err
}

func Consumer(ch *amqp.Channel, name string) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func Declerer(ch amqp.Channel, name string) error {
	_, err := ch.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	return err
}
