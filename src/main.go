package main

import (
	"sync"

	"github.com/dillonmabry/reddit-comments-util/src/batch"
	"github.com/dillonmabry/reddit-comments-util/src/config"
	fileutils "github.com/dillonmabry/reddit-comments-util/src/fileutils"
	graw "github.com/turnage/graw/reddit"
)

func main() {
	// TODO: Setup cmd options for batch processing vs. event sourcing
	// Example of WaitGroup sync, run/process multiple thread comments from channel inputs
	// Sync from channel into common txt file, this assumes non-order of replies
	var wg sync.WaitGroup

	bot := batch.NewBatch(config.BotAgentFile())

	c := make(chan graw.Comment)
	go fileutils.WriteTextOutput("testfile", c)
	wg.Add(2)
	go func() {
		batch.GetAllReplies(bot, "/r/science/comments/6nz1k", c, 0, &wg)
	}()
	go func() {
		batch.GetAllReplies(bot, "/r/AskReddit/comments/3v189r", c, 0, &wg)
	}()
	go func() {
		batch.GetAllReplies(bot, "/r/AskReddit/comments/cliob7", c, 0, &wg)
	}()
	wg.Wait()
	close(c)
}
