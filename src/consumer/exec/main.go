package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/streadway/amqp"

	"github.com/dillonmabry/reddit-processing-utils/src/config"
	"github.com/dillonmabry/reddit-processing-utils/src/datamanager"
	"github.com/dillonmabry/reddit-processing-utils/src/distributed"
	"github.com/dillonmabry/reddit-processing-utils/src/logging"
)

var logger = logging.NewLogger()

// saveMessage persist message to db
func saveMessage(msg []byte) error {
	var postMessage datamanager.PostMessage
	if err := json.Unmarshal(msg, &postMessage); err != nil {
		return err
	}
	saveErr := datamanager.SavePostMessage(&postMessage)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

// handleMessages handler for main goroutine for receiving messages
func handleMessages(msgs <-chan amqp.Delivery, done chan error) {
	for m := range msgs {
		err := saveMessage(m.Body)
		if err != nil {
			logger.Error(err)
			m.Nack(false, false)
		}
		m.Ack(false)
	}
	logger.Info("handleMessages: messages channel closed")
	done <- nil
}

func init() {
	datamanager.InitDB(config.DefaultDb())
}

func main() {
	queue := flag.String("queue", "", "Specifies queue to consume and persist posts")
	flag.Parse()
	c := distributed.NewDistributed(config.DefaultBroker(), *queue)

	msgs, err := c.Channel.Consume(
		c.Queue.Name, // queue string
		"",           // consumer string
		false,        // autoAck bool
		false,        // exclusive bool
		false,        // noLocal bool
		false,        //noWait bool
		nil,          // args amqp.Table
	)
	if err != nil {
		logger.Fatal(fmt.Sprintf("Could not connect to consume for queue %s", c.Queue.Name))
	}

	done := make(chan error)
	handleMessages(msgs, done)
}
