// Package events for reddit comments utils, optional package that can be used for event sourcing
// Intended to use as a standalone package if needed
// Adapted to use event based solution based on https://turnage.gitbooks.io/graw/content/
package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/dillonmabry/reddit-comments-util/src/config"
	distributed "github.com/dillonmabry/reddit-comments-util/src/distributed"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type utilsBot struct {
	bot        reddit.Bot
	mqttClient distributed.Client
}

// NewEvents initialize the graw listener per wrapper
// Based on graw wrapper docs will listen using Go related techniques to check for posts of a subreddit
// Example: events.Init([]string{"bottesting", "science"}, "remind me")
func NewEvents(botAgentFile string, subreddits []string, searchText string) {
	mqttClient := distributed.NewDistributed(config.MQTTBroker(), "topic/test")
	bot, err := reddit.NewBotFromAgentFile(botAgentFile, 0)
	if err != nil {
		panic(err)
	} else {
		cfg := graw.Config{Subreddits: subreddits}
		handler := &utilsBot{bot: bot, mqttClient: mqttClient}
		fmt.Println(fmt.Sprintf("Started bot handler for subreddits: %v", cfg.Subreddits))
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			panic(err)
		} else {
			fmt.Println("graw run failed: ", wait())
		}
	}
}

// Implement the interface per graw
// Listens on Posts per defined subreddit via graw config
func (r *utilsBot) Post(p *reddit.Post) error {
	if strings.Contains(p.SelfText, "remind me now") {
		<-time.After(2 * time.Second) // Buffer
		r.mqttClient.Client.Publish(r.mqttClient.Topic, 0, true, []byte(p.Title))
	}
	return nil
}
