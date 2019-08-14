package main

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/dillonmabry/reddit-comments-util/src/batch"
	"github.com/dillonmabry/reddit-comments-util/src/config"
	fileutils "github.com/dillonmabry/reddit-comments-util/src/fileutils"
	"github.com/turnage/graw/reddit"
	"github.com/urfave/cli"
)

// exportCommentsTxt exports multiple threads replies into single .txt
// subreddit: subreddit, threads: threads to export together
func exportCommentsTxt(subreddit string, threads []string) {
	bot := batch.NewBatch(config.BotAgentFile())

	var wg sync.WaitGroup
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
	bot := batch.NewBatch(config.BotAgentFile())

	var wg sync.WaitGroup
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
		cli.StringFlag{
			Name:  "searchRegex",
			Value: "",
		},
		cli.StringFlag{
			Name:  "csvHeaders",
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
				exportCommentsCsv(c.String("subreddit"), threads,
					[]string{"Accepted", "Application Date", "Decision Date", "Education", "Test Scores", "Experience", "Recommendations", "Comments"},
					`Status: (.*)\n\nApplication Date: (.*)\n\nDecision Date: (.*)\n\nEducation: (.*)\n\nTest Scores: (.*)\n\nExperience: (.*)\n\nRecommendations: (.*)\n\nComments: (.*)`)
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
