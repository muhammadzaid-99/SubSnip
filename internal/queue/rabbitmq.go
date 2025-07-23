// internal/queue/queue.go
package queue

import (
	"log"
	"os"

	"github.com/muhammadzaid-99/SubSnip/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel
var queueName string = "tasks"

func init() {
	var err error
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}
	conn, err = amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	log.Println("RabbitMQ Connected!")

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open channel: %v", err)
	}
	log.Println("RabbitMQ Channel Opened!")

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}
	log.Println("RabbitMQ Queue Declared!")
}

func GetChannel() *amqp.Channel {
	return ch
}

func GetQueueName() string {
	return queueName
}

func Publish(task models.TaskRequest) error {
	body, err := task.ToJSON()
	if err != nil {
		return err
	}
	return ch.Publish(
		"", "tasks", false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func Close() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
