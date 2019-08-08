// Package distributed interactions with Mosquitto MQTT message broker services
// TODO: Setup interactions with mosquitto MQTT with events sourcing
package distributed

import (
	"github.com/dillonmabry/reddit-comments-util/src/logging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Client general client for MQTT
type Client struct {
	Client MQTT.Client
	Topic  string
}

// NewDistributed create a distributed client to publish events to broker
// brokerURL: broker to connect, topic: topic to publish to
func NewDistributed(brokerURL string, topic string, msgHandler MQTT.MessageHandler) Client {
	logger := logging.NewLogger()
	opts := MQTT.NewClientOptions().AddBroker(brokerURL)
	//opts.SetClientID("mac-go")
	opts.SetDefaultPublishHandler(msgHandler)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Fatal("Could not connect publishing client")
	} else {
		logger.Info("Connected to broker service")
	}
	return Client{Client: client, Topic: topic}
}
