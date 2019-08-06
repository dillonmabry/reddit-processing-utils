// Package config options
package config

import (
	"os"
)

// BotAgentFile to return
func BotAgentFile() string {
	botFileName := os.Getenv("GRAW_BOT_AGENT")
	if botFileName == "" {
		panic("botFileName Env variable not found")
	}
	return botFileName
}

//TODO: Add MQTT broker config
