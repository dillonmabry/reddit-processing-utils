// Package batch any batch level manipulations via graw
package batch

import (
	"regexp"
	"sync"

	"github.com/dillonmabry/reddit-comments-util/src/logging"
	"github.com/turnage/graw/reddit"
	gwrap "github.com/turnage/graw/reddit"
)

// NewBatch creates batch style handler for graw bot interactions
// botAgentFile: the name of the bot agent file to use
func NewBatch(botAgentFile string) gwrap.Bot {
	logger := logging.NewLogger()
	bot, err := reddit.NewBotFromAgentFile(botAgentFile, 0)
	if err != nil {
		logger.Error("Could not create bot agent from file")
		panic(err)
	}
	logger.Info("Started bot for batch processing")
	return bot
}

// GetAllReplies to get all comments of a particular thread
// bot: reddit Bot, thread: thread, c: channel, wg WaitGroup ref
func GetAllReplies(bot gwrap.Bot, thread string, c chan gwrap.Comment, wg *sync.WaitGroup) {
	harvest, err := bot.Thread(thread)
	if err != nil {
		panic(err)
	}

	for _, comment := range harvest.Replies {
		for _, comment := range comment.Replies {
			c <- *comment
		}
		c <- *comment
	}
	defer wg.Done()
}

// findMatches filters via regex into designated channel
// regex: reference to regex compiler, text: text to match, c: channel
func findMatches(regex *regexp.Regexp, text string, c chan []string) {
	matches := regex.FindAllStringSubmatch(text, -1)
	for m := range matches {
		//fmt.Println([]string{matches[m][1]})
		c <- []string{matches[m][1]}
	}
}

// GetFilteredReplies to get all comments of a particular thread
// bot: reddit Bot, thread: thread, c: channel, searchPattern regexp, wg WaitGroup ref
func GetFilteredReplies(bot gwrap.Bot, thread string, c chan []string, searchPattern string, wg *sync.WaitGroup) {
	harvest, err := bot.Thread(thread)
	if err != nil {
		panic(err)
	}
	r, err := regexp.Compile("(?i)" + searchPattern)
	if err != nil {
		panic(err)
	}

	for _, comment := range harvest.Replies {
		for _, comment := range comment.Replies {
			findMatches(r, comment.Body, c)
		}
		findMatches(r, comment.Body, c)
	}
	defer wg.Done()
}
