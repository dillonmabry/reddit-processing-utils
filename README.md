# Reddit Processing Utils

Wrapper for Reddit API including batch processing and events/broker services

## Goals
- Create scalable services that can listen on multiple subreddits and gather data concurrently
- Also allow fast batch processing of multiple threads, gathering data as needed
- Setup broker interface to connect with and allow multiple topics to be processed

## Install Instructions

## Run Instructions
```
export GRAW_BOT_AGENT=localbot.agent && GRAW_BOT_SUBREDDITS=bottesting,science && time go run src/main.go events --subreddits AskReddit,science --searchText <your text you wish to search>
```

## Tests
