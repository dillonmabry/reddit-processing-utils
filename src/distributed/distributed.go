// Package distributed interactions with Mosquitto MQTT message broker services
// TODO: Setup interactions with mosquitto MQTT with events sourcing
package distributed

import (
	"fmt"
	"log"

	"github.com/dillonmabry/reddit-comments-util/src/logging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Client general client for MQTT
type Client struct {
	Client MQTT.Client
	Topic  string
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	log.Println(fmt.Sprintf("TOPIC: %s\n", msg.Topic()))
	log.Println(fmt.Sprintf("MSG: %s\n", msg.Payload()))
}

// NewDistributed create a distributed client to publish events to broker
// brokerURL: broker to connect, topic: topic to publish to
func NewDistributed(brokerURL string, topic string) Client {
	logger := logging.NewLogger()
	opts := MQTT.NewClientOptions().AddBroker(brokerURL)
	opts.SetClientID("mac-go")
	opts.SetDefaultPublishHandler(f)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logger.Error(token.Error())
		panic(token.Error())
	} else {
		logger.Info("Connected to broker service")
	}
	return Client{Client: client, Topic: topic}
}
