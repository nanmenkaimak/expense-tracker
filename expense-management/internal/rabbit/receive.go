package rabbit

import (
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func ReceiveMessage(username string) (string, error) {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		return "", errors.Wrap(err, "rabbit connect")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return "", errors.Wrap(err, "rabbit channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("token_%s", username),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", errors.Wrap(err, "rabbit queue declare")
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", errors.Wrap(err, "rabbitmq declare consumer")
	}
	data := <-msgs

	err = data.Ack(false)
	if err != nil {
		return "", err
	}

	return string(data.Body), nil
}
