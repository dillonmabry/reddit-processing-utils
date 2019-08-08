package consumer

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type consumer struct {
	msgHandler MQTT.MessageHandler
}

// NewConsumer create new consumer with message handler
func NewConsumer(msgHandler MQTT.MessageHandler) *consumer {
	consumer := consumer{msgHandler: msgHandler}
	return &consumer
}
