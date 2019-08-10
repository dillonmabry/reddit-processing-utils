package main

import (
	"flag"
	"os"
	"strings"

	"github.com/dillonmabry/reddit-comments-util/src/config"
	"github.com/dillonmabry/reddit-comments-util/src/events"
)

func main() {
	subreddits := flag.String("subreddits", "", "List of comma separated subreddits to listen on")
	searchText := flag.String("search", "", "Text to search for inside of post body for listening")
	queue := flag.String("queue", "", "Main publishing queue, to be used in coordination with consumer queue")
	flag.Parse()

	subredditsList := strings.Split(*subreddits, ",")
	events.NewEvents(config.BotAgentFile(), *queue, subredditsList, *searchText)
	defer os.Exit(0)
}
