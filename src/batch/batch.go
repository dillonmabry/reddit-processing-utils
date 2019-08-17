// Package batch any batch level manipulations via graw
package batch

import (
	"fmt"
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
// regex: reference to regex compiler, groupCount int number of capture groups, text: text to match, c: channel
func findMatches(regex *regexp.Regexp, groupCount int, text string, c chan []string) {
	matches := regex.FindAllStringSubmatch(text, -1)
	for m := range matches {
		var groupMatches []string
		for i := 1; i < groupCount; i++ {
			groupMatches = append(groupMatches, matches[m][i])
		}
		c <- groupMatches
	}
}

// getFilteredReplies to get all comments of a particular thread
// bot: reddit Bot, thread: thread, c: channel, regex regexp
func getFilteredReplies(bot gwrap.Bot, thread string, c chan []string, regex string) {
	harvest, err := bot.Thread(thread)
	if err != nil {
		panic(err)
	}
	r, err := regexp.Compile("(?i)" + regex)
	if err != nil {
		panic(err)
	}

	for _, comment := range harvest.Replies {
		for _, comment := range comment.Replies {
			findMatches(r, len(r.SubexpNames()), comment.Body, c)
		}
		findMatches(r, len(r.SubexpNames()), comment.Body, c)
	}
}

// RepliesProducer creates a goroutine to return a receiving channel of reddit replies to a thread
// bot: reddit bot, subreddit: subreddit, thread: thread, regex: regex search pattern form based
func RepliesProducer(bot gwrap.Bot, subreddit string, thread string, regex string) <-chan []string {
	out := make(chan []string)
	go func() {
		getFilteredReplies(bot, fmt.Sprintf("/r/%s/comments/%s", subreddit, thread), out, regex)
		close(out)
	}()
	return out
}

// RepliesProducerMerged creates a goroutine to return a receiving channel of reddit replies to a thread
// Uses a shared channel to merge results
// bot: reddit bot, subreddit: subreddit, thread: thread, regex: regex search pattern form based
func RepliesProducerMerged(bot gwrap.Bot, subreddit string, thread string, regex string, c chan []string) <-chan []string {
	go func(thread string) {
		getFilteredReplies(bot, fmt.Sprintf("/r/%s/comments/%s", subreddit, thread), c, regex)
		close(c)
	}(thread)
	return c
}
