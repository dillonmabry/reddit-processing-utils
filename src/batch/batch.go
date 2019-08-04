// Package batch any batch level manipulations via graw
package batch

import (
	"sync"

	"github.com/turnage/graw/reddit"
	graw "github.com/turnage/graw/reddit"
)

// GetAllReplies to get all comments of a particular thread
// bot: reddit Bot, thread: thread, c: channel, top: max posts
func GetAllReplies(bot graw.Bot, thread string, c chan graw.Comment, top int, wg *sync.WaitGroup) {
	harvest, err := bot.Thread(thread)
	if err != nil {
		panic(err)
	}

	if top != 0 {
		for _, comment := range harvest.Replies[:top] {
			c <- *comment
		}
	}

	// Separating these into goroutines does not help performance
	for _, comment := range harvest.Replies {
		for _, comment := range comment.Replies {
			c <- *comment
		}
		c <- *comment
	}
	defer wg.Done()
}

// NewBatch creates batch style handler for graw bot interactions
// botAgentFile: the name of the bot agent file to use
func NewBatch(botAgentFile string) graw.Bot {
	bot, err := reddit.NewBotFromAgentFile(botAgentFile, 0)
	if err != nil {
		panic(err)
	}
	return bot
}
