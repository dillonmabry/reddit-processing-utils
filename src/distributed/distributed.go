// Package distributed interactions with Mosquitto MQTT message broker services
package distributed

import (
	"fmt"

	"github.com/dillonmabry/reddit-comments-util/src/logging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Client general for MQTT
type Client struct {
	Client MQTT.Client
	Topic  string
}

// NewDistributed create a distributed client to publish events to broker
// brokerURL: broker to connect, topic: topic to publish to, handler: function handler for publishing
func NewDistributed(brokerURL string, topic string, handler MQTT.MessageHandler) *Client {
	logger := logging.NewLogger()
	opts := MQTT.NewClientOptions().AddBroker(brokerURL)
	opts.SetDefaultPublishHandler(handler)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Fatal("Could not connect publishing client")
		return nil
	}
	logger.Info(fmt.Sprintf("Connected to broker service via topic: %s", topic))
	distClient := Client{Client: client, Topic: topic}
	return &distClient
}
