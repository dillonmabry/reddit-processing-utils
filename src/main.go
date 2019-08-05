package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dillonmabry/reddit-comments-util/src/batch"
	"github.com/dillonmabry/reddit-comments-util/src/config"
	"github.com/dillonmabry/reddit-comments-util/src/events"
	"github.com/dillonmabry/reddit-comments-util/src/fileutils"
	graw "github.com/turnage/graw/reddit"
	"github.com/urfave/cli"
)

// exportCommentsTxt exports multiple threads replies into single .txt
// subreddit: subreddit, threads: threads to export together
func exportCommentsTxt(subreddit string, threads []string) {
	var wg sync.WaitGroup
	bot := batch.NewBatch(config.BotAgentFile())

	c := make(chan graw.Comment)
	go fileutils.WriteTextOutput("testfile", c)
	wg.Add(len(threads))
	for _, thread := range threads {
		go func(thread string) {
			batch.GetAllReplies(bot, fmt.Sprintf("/r/%s/comments/%s", subreddit, thread), c, &wg)
		}(thread)
	}
	wg.Wait()
	close(c)
}

// exportCommentsTxt exports multiple threads replies into single .txt
// subreddit: subreddit, threads: threads to export together, searchPattern regexp pattern
func exportCommentsCsv(subreddit string, threads []string, headers []string, searchPattern string) {
	var wg sync.WaitGroup
	bot := batch.NewBatch(config.BotAgentFile())

	c := make(chan []string)
	go fileutils.WriteCsvOutput("testfile", headers, c)
	wg.Add(len(threads))
	for _, thread := range threads {
		go func(thread string) {
			batch.GetFilteredReplies(bot, fmt.Sprintf("/r/%s/comments/%s", subreddit, thread), c, searchPattern, &wg)
		}(thread)
	}
	wg.Wait()
	close(c)
}

// subredditsListener listens to multiple subreddits for specific post content
// subreddits: subreddits, searchText: text to listen for inside a specific post body
func subRedditsListener(subreddits []string, searchText string) {
	events.NewEvents(config.BotAgentFile(), subreddits, searchText)
}

func main() {
	app := cli.NewApp()
	app.Name = "Reddit Comments Utility CLI"
	app.Usage = "Allows batch export of reddit comments and event listening for subreddits"

	var flags = []cli.Flag{
		cli.StringFlag{
			Name:  "subreddit",
			Value: "",
		},
		cli.StringFlag{
			Name:  "threads",
			Value: "",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "batch",
			Usage: "Export multiple reddit threads to single .txt file",
			Flags: flags,
			Action: func(c *cli.Context) error {
				threads := strings.Split(c.String("threads"), ",")
				exportCommentsTxt(c.String("subreddit"), threads)
				return nil
			},
		},
		// {
		// 	Name:  "event",
		// 	Usage: "Listen to specified search string from multiple subreddits then notify",
		// 	Flags: flags,
		// 	Action: func(c *cli.Context) error {
		// 		subreddits := strings.Split(c.String("subreddits"), ",")
		// 		subRedditsListener(subreddits, c.String("searchText"))
		// 		return nil
		// 	},
		// },
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
