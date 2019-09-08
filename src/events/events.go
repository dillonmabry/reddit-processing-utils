// Package events for reddit comments utils, optional package that can be used for event sourcing
// Intended to use as a standalone package if needed
// Adapted to use event based solution based on https://turnage.gitbooks.io/graw/content/
package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dillonmabry/reddit-processing-utils/src/config"
	"github.com/dillonmabry/reddit-processing-utils/src/datamanager"
	"github.com/dillonmabry/reddit-processing-utils/src/distributed"
	"github.com/dillonmabry/reddit-processing-utils/src/logging"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var logger = logging.NewLogger()

// searchBot bot for events sourcing
type searchBot struct {
	bot        reddit.Bot
	distClient *distributed.Client
	searchText string
}

// NewEvents initialize the graw listener per wrapper
// Based on graw wrapper docs will listen using Go related techniques to check for posts of a subreddit
// botAgentFile: bot agent local, queueName: queue of events publisher, subreddits: subreddits, searchText: text contains
func NewEvents(botAgentFile string, queueName string, subreddits []string, searchText string) {

	distClient := distributed.NewDistributed(config.DefaultBroker(), queueName)

	bot, err := reddit.NewBotFromAgentFile(botAgentFile, 0)
	if err != nil {
		logger.Fatal("Could not create bot agent from file")
	} else {
		cfg := graw.Config{Subreddits: subreddits}
		handler := &searchBot{bot: bot, distClient: distClient, searchText: searchText}
		logger.Info(fmt.Sprintf("Started bot handler for subreddits: %v", strings.Join(cfg.Subreddits, ",")))
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			logger.Fatal("Graw run failed from handler setup")
		} else {
			logger.Info("graw run failed: ", wait())
		}
	}
}

// Implement the interface per graw
// Listens on Posts per defined subreddit via graw config
func (r *searchBot) Post(p *reddit.Post) error {
	if strings.Contains(p.SelfText, r.searchText) {
		<-time.After(2 * time.Second) // Buffer
		post := datamanager.PostMessage{URL: p.URL, Text: p.SelfText}
		msg, err := json.Marshal(post)
		if err != nil {
			logger.Error(fmt.Sprintf("Error converting to JSON for Reddit post %s", p.URL))
		}
		pubErr := r.distClient.Channel.Publish(
			config.DefaultExchange(),
			r.distClient.Queue.Name,
			false,
			false,
			distributed.PublishBody(msg),
		)
		if pubErr != nil {
			logger.Error(fmt.Sprintf("Error publishing message to queue %s", r.distClient.Queue.Name))
		}
	}
	return nil
}
