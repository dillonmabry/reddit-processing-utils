// Package distributed interactions with Mosquitto MQTT message broker services
// TODO: Setup interactions with mosquitto MQTT with events sourcing
package distributed

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// Client general client for MQTT
type Client struct {
	Client MQTT.Client
	Topic  string
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

// NewDistributed create a distributed client to publish events to broker
// brokerURL: broker to connect
func NewDistributed(brokerURL string, topic string) Client {
	opts := MQTT.NewClientOptions().AddBroker("tcp://" + brokerURL)
	opts.SetClientID("mac-go")
	opts.SetDefaultPublishHandler(f)

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	return Client{Client: client, Topic: topic}
}
