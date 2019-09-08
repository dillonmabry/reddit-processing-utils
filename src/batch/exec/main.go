package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dillonmabry/reddit-processing-utils/src/batch"
	"github.com/dillonmabry/reddit-processing-utils/src/config"
	"github.com/dillonmabry/reddit-processing-utils/src/fileutils"
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

// exportCommentsCsvMerged exports multiple threads replies into single .csv
// subreddit: subreddit, threads: threads to export together, regex regexp pattern
func exportCommentsCsvMerged(subreddit string, threads []string, headers []string, regex string) {
	bot := batch.NewBatch(config.BotAgentFile())

	fname := strings.Join(threads[:], "_") + ".csv"
	w, err := fileutils.NewCsvWriter(fname)
	if err != nil {
		panic(err)
	}
	w.Write(headers)
	defer w.Flush()

	for _, thread := range threads {
		c := make(chan []string)
		for msg := range batch.RepliesProducerMerged(bot, subreddit, thread, regex, c) {
			w.Write(msg)
		}
	}
}

// exportCommentsCsv exports multiple threads replies into multiple .csv files
// subreddit: subreddit, threads: threads to export together, regex regexp pattern
func exportCommentsCsv(subreddit string, threads []string, headers []string, regex string) {
	bot := batch.NewBatch(config.BotAgentFile())

	for _, thread := range threads {
		fname := thread + ".csv"
		w, err := fileutils.NewCsvWriter(fname)
		if err != nil {
			panic(err)
		}
		w.Write(headers)
		defer w.Flush()
		for msg := range batch.RepliesProducer(bot, subreddit, thread, regex) {
			w.Write(msg)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "Reddit Comments Utility - Batch CLI"
	app.Usage = "Allows batch-type functionality of Reddit API including comments, search, and more"

	var flags = []cli.Flag{
		cli.StringFlag{
			Name:  "subreddit",
			Value: "",
			Usage: "Subreddit group",
		},
		cli.StringFlag{
			Name:  "threads",
			Value: "",
			Usage: "Reddit threads comma separated, ie: thread1,thread2",
		},
		cli.StringFlag{
			Name:  "headers",
			Value: "",
			Usage: `Csv headers which match the "form" style regex search`,
		},
		cli.BoolFlag{
			Name:  "m",
			Usage: "Indicates whether to merge all csv threads into a single file",
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
				subreddit := c.String("subreddit")
				threads := strings.Split(c.String("threads"), ",")
				headers := strings.Split(c.String("headers"), ",")
				searchPattern := fileutils.HeadersToRegex(headers)
				isMerge := c.Bool("m")
				if isMerge {
					exportCommentsCsvMerged(subreddit, threads, headers, searchPattern)
					return nil
				}
				exportCommentsCsv(subreddit, threads, headers, searchPattern)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
