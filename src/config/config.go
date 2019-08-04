// Package config options
package config

import (
	"os"
	"strings"
)

// BotAgentFile to return
func BotAgentFile() string {
	envVar := "GRAW_BOT_AGENT"
	botFileName := os.Getenv(envVar)
	if botFileName == "" {
		panic("Env variable not found")
	}
	return botFileName
}

// SubReddits to return, must be comma separated string
func SubReddits() []string {
	envVar := "GRAW_BOT_SUBREDDITS"
	subreddits := os.Getenv(envVar)
	if subreddits == "" {
		panic("Env variable not found")
	}
	return strings.Split(subreddits, ",")
}
