// Package config options
package config

import (
	"os"
)

// BotAgentFile to return
func BotAgentFile() string {
	envVar := "GRAW_BOT_AGENT"
	botFileName := os.Getenv(envVar)
	if botFileName == "" {
		panic("botFileName Env variable not found")
	}
	return botFileName
}

// MQTTBroker to return
func MQTTBroker() string {
	envVar := "MQTT_BROKER"
	MQTTBroker := os.Getenv(envVar)
	if MQTTBroker == "" {
		panic("MQTTBroker Env variable not found")
	}
	return MQTTBroker
}
