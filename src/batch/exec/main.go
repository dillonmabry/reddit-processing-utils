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
	bot := batch.NewBatch(config.BotAgentFile())

	var wg sync.WaitGroup
	c := make(chan reddit.Comment)
	fname := strings.Join(threads[:], "_") + ".txt"
	go fileutils.WriteTextOutput(fname, c)
	wg.Add(len(threads))
	for _, thread := range threads {
		go func(thread string) {
			batch.GetAllReplies(bot, fmt.Sprintf("/r/%s/comments/%s", subreddit, thread), c, &wg)
		}(thread)
	}
	wg.Wait()
	close(c)
}

// exportCommentsTxt exports multiple threads replies into single .csv
// subreddit: subreddit, threads: threads to export together, headers: csv headers, searchPattern regexp pattern
func exportCommentsCsv(subreddit string, threads []string, headers []string, searchPattern string) {
	bot := batch.NewBatch(config.BotAgentFile())

	var wg sync.WaitGroup
	c := make(chan []string)
	wg.Add(len(threads))

	fname := strings.Join(threads[:], "_") + ".csv"
	w, err := fileutils.NewCsvWriter(fname)
	w.Write(headers)
	defer w.Flush()

	if err != nil {
		panic(err)
	}
	for _, thread := range threads {
		go func(thread string) {
			batch.GetFilteredReplies(bot, fmt.Sprintf("/r/%s/comments/%s", subreddit, thread), c, searchPattern, &wg)
		}(thread)
	}
	go func() { //For some reason this is blocking when not in a separate goroutine
		for msg := range c {
			w.Write(msg)
		}
	}()
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
		cli.StringFlag{
			Name:  "headers",
			Value: "",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "txt",
			Usage: "Export multiple reddit threads to single .txt file",
			Flags: flags,
			Action: func(c *cli.Context) error {
				threads := strings.Split(c.String("threads"), ",")
				exportCommentsTxt(c.String("subreddit"), threads)
				return nil
			},
		},
		{
			Name:  "csv",
			Usage: "Export regex based search criteria to csv or multiple csvs for reddit threads",
			Flags: flags,
			Action: func(c *cli.Context) error {
				threads := strings.Split(c.String("threads"), ",")
				headers := strings.Split(c.String("headers"), ",")
				searchPattern := fileutils.HeadersToRegex(headers)
				exportCommentsCsv(c.String("subreddit"), threads, headers, searchPattern)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
