// Package events for reddit comments utils, optional package that can be used for event sourcing
// Intended to use as a standalone package if needed
// Adapted to use event based solution based on https://turnage.gitbooks.io/graw/content/
package events

import (
	"fmt"
	"strings"
	"time"

	distributed "github.com/dillonmabry/reddit-comments-util/src/distributed"
	logging "github.com/dillonmabry/reddit-comments-util/src/logging"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type searchBot struct {
	bot        reddit.Bot
	mqttClient distributed.Client
	searchText string
}

// NewEvents initialize the graw listener per wrapper
// Based on graw wrapper docs will listen using Go related techniques to check for posts of a subreddit
// botAgentFile: bot agent local, subreddits: subreddits, searchText: text contains
func NewEvents(botAgentFile string, subreddits []string, searchText string) {
	logger := logging.NewLogger()
	mqttClient := distributed.NewDistributed("tcp://192.168.1.220:1883", "topic/test")
	bot, err := reddit.NewBotFromAgentFile(botAgentFile, 0)

	if err != nil {
		logger.Error(err)
		panic(err)
	} else {
		cfg := graw.Config{Subreddits: subreddits}
		handler := &searchBot{bot: bot, mqttClient: mqttClient, searchText: searchText}
		logger.Info(fmt.Sprintf("Started bot handler for subreddits: %v", cfg.Subreddits))
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			logger.Error(err)
			panic(err)
		} else {
			logger.Info("graw run failed: ", wait())
		}
	}
}

// Implement the interface per graw
// Listens on Posts per defined subreddit via graw config
func (r *searchBot) Post(p *reddit.Post) error {
	if strings.Contains(p.SelfText, r.searchText) {
		<-time.After(2 * time.Second)                                                // Buffer
		r.mqttClient.Client.Publish(r.mqttClient.Topic, 0, true, []byte(p.SelfText)) // Publish source URL
	}
	return nil
}
