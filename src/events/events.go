// Package events for reddit comments utils, optional package that can be used for event sourcing
// Intended to use as a standalone package if needed
// Adapted to use event based solution based on https://turnage.gitbooks.io/graw/content/
package events

import (
	"fmt"
	"strings"
	"time"

	"github.com/dillonmabry/reddit-comments-util/src/config"
	"github.com/dillonmabry/reddit-comments-util/src/logwrapper"
	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"
)

type utilsBot struct {
	bot reddit.Bot
}

var reminderText string

// Implement the interface per graw
// Listens on Posts per defined subreddit via graw config
func (r *utilsBot) Post(p *reddit.Post) error {
	if strings.Contains(p.SelfText, reminderText) {
		<-time.After(2 * time.Second) // Buffer
		return r.bot.SendMessage(
			p.Author,
			fmt.Sprintf("Reminder: %s", p.Title),
			"You've been reminded!",
		)
	}
	return nil
}

// Init initialize the graw listener per wrapper, this logic is wrapped in our Init func
// Based on graw wrapper docs will listen using Go related techniques to check for posts of a subreddit
// Example: events.Init([]string{"bottesting", "science"}, "remind me")
func Init(subreddits []string, searchText string) {
	logger := logwrapper.NewLogger()
	reminderText = searchText
	bot, err := reddit.NewBotFromAgentFile(config.BotAgentFile(), 0)
	if err != nil {
		logger.Error("Failed to create bot handle: ", err)
		panic(err)
	} else {
		cfg := graw.Config{Subreddits: subreddits}
		handler := &utilsBot{bot: bot}
		logger.Info(fmt.Sprintf("Started bot handler for subreddits: %v", cfg.Subreddits))
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			logger.Error("Failed to start graw run: ", err)
			panic(err)
		} else {
			logger.Error("graw run failed: ", wait())
		}
	}
}
