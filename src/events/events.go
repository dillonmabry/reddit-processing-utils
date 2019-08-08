// Package events for reddit comments utils, optional package that can be used for event sourcing
// Intended to use as a standalone package if needed
// Adapted to use event based solution based on https://turnage.gitbooks.io/graw/content/
package events

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dillonmabry/reddit-comments-util/src/datamanager"
	"github.com/dillonmabry/reddit-comments-util/src/distributed"
	"github.com/dillonmabry/reddit-comments-util/src/logging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

var logger = logging.NewLogger()

// searchBot bot for events sourcing
type searchBot struct {
	bot        reddit.Bot
	mqttClient *distributed.Client
	searchText string
}

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	logger.Info(fmt.Sprintf("TOPIC: %s\n", msg.Topic()))
	logger.Info(fmt.Sprintf("MSG: %s\n", msg.Payload()))
}

// NewEvents initialize the graw listener per wrapper
// Based on graw wrapper docs will listen using Go related techniques to check for posts of a subreddit
// botAgentFile: bot agent local, subreddits: subreddits, searchText: text contains
func NewEvents(botAgentFile string, subreddits []string, searchText string) {
	//TODO: Make topic selection generic
	mqttClient := distributed.NewDistributed("tcp://192.168.1.220:1883", "topic/test", f)

	bot, err := reddit.NewBotFromAgentFile(botAgentFile, 0)
	if err != nil {
		logger.Fatal("Could not create bot agent from file")
	} else {
		cfg := graw.Config{Subreddits: subreddits}
		handler := &searchBot{bot: bot, mqttClient: mqttClient, searchText: searchText}
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
		message := datamanager.PostMessage{URL: p.URL, Text: p.SelfText}
		messageJSON, err := json.Marshal(message)
		if err != nil {
			logger.Error(fmt.Sprintf("Error converting to JSON for Reddit post %s", p.URL))
		}
		if token := r.mqttClient.Client.Publish(r.mqttClient.Topic, 0, false, messageJSON); token.Wait() && token.Error() != nil {
			logger.Fatal(token.Error())
		}
	}
	return nil
}
