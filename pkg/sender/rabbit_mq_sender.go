package sender

import (
	"fmt"
	"log"

	"receiving_service/internal/metrics"

	"github.com/streadway/amqp"
)

type Rabbit_Handler struct {
}

func (r Rabbit_Handler) StartServer(addr string, tr *chan []byte) {
	metrics.PromoHandler{}.StartMetrics(addr)
	fmt.Println("RabbitMq!")

	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ Instance")

	file := <-*tr

	mes, err := r.Upload(file, "client", *conn)

	if err != nil {
		log.Fatal(err)
	}

	*tr <- mes

}

func (r Rabbit_Handler) Upload(file []byte, name string, conn amqp.Connection) ([]byte, error) {
	metrics.PromoHandler{}.RecordMetrics()
	ch, err := conn.Channel()

	if err != nil {
		log.Fatal(err)
	}
	defer ch.Close()

	err = Publisher(*ch, "Server", []byte(name))

	if err != nil {
		return []byte{}, err
	}

	be := 0
	en := 1024
	for {

		if len(file) < en {
			err = Publisher(*ch, name, file[be:])
			if err != nil {
				return []byte{}, err
			}

			err = Publisher(*ch, name, []byte("Stop"))
			if err != nil {
				return []byte{}, err
			}

			break
		}
		err = Publisher(*ch, name, file[be:en])

		if err != nil {
			return []byte{}, err
		}

		be = en
		en += 1024
	}

	fmt.Println("Successfully Published Message to Queue")
	ch, err = conn.Channel()

	if err != nil {
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

func Publisher(ch amqp.Channel, name string, mes []byte) error {
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