package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dillonmabry/reddit-comments-util/src/batch"
	"github.com/dillonmabry/reddit-comments-util/src/config"
	"github.com/dillonmabry/reddit-comments-util/src/fileutils"
	"github.com/turnage/graw/reddit"
	"github.com/urfave/cli"
)

// exportCommentsTxt exports multiple threads replies into single .txt
// subreddit: subreddit, threads: threads to export together
func exportCommentsTxt(subreddit string, threads []string) {
	var wg sync.WaitGroup
	bot := batch.NewBatch(config.BotAgentFile())

	c := make(chan reddit.Comment)
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

func main() {
	app := cli.NewApp()
	app.Name = "Reddit Comments Utility - Batch CLI"
	app.Usage = "Allows batch-type functionality of Reddit API including comments, search, and more"

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
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
