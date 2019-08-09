# Reddit Processing Utils

Wrapper for Reddit API including batch processing and events/broker services

## Goals
- Create scalable services that can listen on multiple subreddits and gather data concurrently
- Also allow fast batch processing of multiple threads, gathering data as needed
- Setup broker interface to connect with and allow multiple topics to be processed

## Install Instructions

## Run Instructions
Modify/run .bat files as needed
OR
Run manually:

Event listener
```
export GRAW_BOT_AGENT=localbot.agent && go run src/events/exec/main.go events --subreddits AskReddit,science --searchText <your text you wish to search>
```

Consumer
```
go run src/consumer/exec/main.go
```
## Tests
