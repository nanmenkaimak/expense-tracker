package rabbit

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"time"
)

func SendMessage(data []byte, id string) error {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ"))
	if err != nil {
		return errors.Wrap(err, "rabbit connect")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return errors.Wrap(err, "rabbit channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		fmt.Sprintf("token_%s", id),
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.Wrap(err, "rabbit queue declare")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = ch.PublishWithContext(
		ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
	if err != nil {
		return errors.Wrap(err, "rabbit publish message")
	}
	return nil
}
