// Package distributed interactions with AMQP RabbitMQ
package distributed

import (
	"fmt"

	"github.com/dillonmabry/reddit-processing-utils/src/config"
	"github.com/dillonmabry/reddit-processing-utils/src/logging"
	"github.com/streadway/amqp"
)

// Client general for AMQP
type Client struct {
	Channel *amqp.Channel
	Queue   amqp.Queue
}

// PublishBody return publish object
func PublishBody(data []byte) amqp.Publishing {
	return amqp.Publishing{Body: data}
}

// NewDistributed create a distributed client for general use
// brokerURL: broker to connect, queueName: queue to publish to
func NewDistributed(brokerURL string, queueName string) *Client {
	logger := logging.NewLogger()

	conn, err := amqp.Dial(brokerURL)
	if err != nil {
		logger.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("Failed to open channel")
	}

	if err := ch.ExchangeDeclare(
		config.DefaultExchange(), // name of the exchange
		"fanout",                 // type
		true,                     // durable
		false,                    // delete when complete
		false,                    // internal
		false,                    // noWait
		nil,                      // arguments
	); err != nil {
		logger.Fatal(fmt.Sprintf("Exchange: %s", err))
	}

	q, err := ch.QueueDeclare(
		queueName, // name of the queue
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args table
	)
	if err := ch.QueueBind(
		q.Name,                   // name of the queue
		"",                       // binding key
		config.DefaultExchange(), // source exchange
		false,                    // noWait
		nil,                      // arguments
	); err != nil {
		logger.Fatal(fmt.Sprintf("Queue Bind: %s", err))
	}

	distClient := Client{Channel: ch, Queue: q}
	logger.Info(fmt.Sprintf("Connected to broker service via queue: %s", q.Name))
	return &distClient
}
