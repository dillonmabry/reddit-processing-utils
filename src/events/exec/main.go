package main

import (
	"os"
	"strings"

	"github.com/dillonmabry/reddit-comments-util/src/config"
	"github.com/dillonmabry/reddit-comments-util/src/events"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Reddit Comments Utility - Events CLI"
	app.Usage = "Allows event-based listening via Reddit with different functions"

	var flags = []cli.Flag{
		cli.StringFlag{
			Name:  "subreddits",
			Value: "",
		},
		cli.StringFlag{
			Name:  "searchText",
			Value: "",
		},
		cli.StringFlag{
			Name:  "topic",
			Value: "",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "events",
			Usage: "Listen to specified search string from multiple subreddits then notify",
			Flags: flags,
			Action: func(c *cli.Context) error {
				subreddits := strings.Split(c.String("subreddits"), ",")
				events.NewEvents(config.BotAgentFile(), c.String("topic"), subreddits, c.String("searchText"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
