// Package config options
package config

import (
	"os"
)

// README: Edit your default exchanges if needed
const defaultBotFile = "localbot.agent"
const defaultBroker = "amqp://admin:admin@192.168.44.128:5672"
const defaultExchange = "amq.fanout"
const defaultDbConn = "host=localhost port=5432 user=admin password=admin dbname=reddit sslmode=disable"

// BotAgentFile to return
func BotAgentFile() string {
	botFileName := os.Getenv("GRAW_BOT_AGENT")
	if botFileName == "" {
		return defaultBotFile
	}
	return botFileName
}

// DefaultBroker default broker url to connect
func DefaultBroker() string {
	brokerName := os.Getenv("AMQP_DEFAULT_BROKER")
	if brokerName == "" {
		return defaultBroker
	}
	return brokerName
}

// DefaultExchange default amqp exchange for project
func DefaultExchange() string {
	exchangeName := os.Getenv("AMQP_DEFAULT_EXCHANGE")
	if exchangeName == "" {
		return defaultExchange
	}
	return exchangeName
}

//DefaultDb default database for Postgres
func DefaultDb() string {
	dbConn := os.Getenv("POSTGRES_CONN")
	if dbConn == "" {
		return defaultDbConn
	}
	return dbConn
}
