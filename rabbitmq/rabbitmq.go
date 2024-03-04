package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	QueueName = "Notifications"
)

var conn *amqp091.Connection
var channel *amqp091.Channel

func InitRabbitMq() error {
	var err error
	conn, err := amqp091.Dial("amqp091://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ :%v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to return channel :%v", err)
	}

	_, err = channel.QueueDeclare(
		"Notifications",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %v", err)
	}
	return nil
}

func Close() {
	if channel != nil {
		channel.Close()
	}

	if conn != nil {
		conn.Close()
	}
}

func PublishNotification(message string, userID int) error {
	notification := map[string]interface{}{
		"message": message,
		"userID":  userID,
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to Marshal notification: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = channel.PublishWithContext(
		ctx,
		"",
		"TestNotification",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body: body,
		},
	)
	if err != nil {
		return fmt.Errorf("could not publish : %v", err)
	}
	return nil
}
