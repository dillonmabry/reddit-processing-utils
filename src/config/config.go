// Package config options
package config

import (
	"os"
)

// const: TODO: Add as environment variables to configure
const defaultBroker = "amqp://admin:admin@192.168.1.113:5672"
const defaultExchange = "amq.fanout"

// BotAgentFile to return
func BotAgentFile() string {
	botFileName := os.Getenv("GRAW_BOT_AGENT")
	if botFileName == "" {
		panic("botFileName Env variable not found")
	}
	return botFileName
}

// DefaultExchange default amqp exchange for project
func DefaultExchange() string {
	return defaultExchange
}

// DefaultBroker default broker url to connect
func DefaultBroker() string {
	return defaultBroker
}
