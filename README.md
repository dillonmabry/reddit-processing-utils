# Reddit Processing Utils
[![Build Status](https://goreportcard.com/badge/github.com/dillonmabry/reddit-processing-utils)](https://goreportcard.com/report/github.com/dillonmabry/reddit-processing-utils)

Wrapper for Reddit API including batch processing and events/broker services

## Goals
- Create scalable services that can listen on multiple subreddits and gather data concurrently
- Also allow fast batch processing of multiple threads, gathering data as needed
- Setup broker interface to connect with and allow multiple topics to be processed

## Install Instructions
1. Get:
```
go get -u github.com/dillonmabry/reddit-processing-utils/...
```
2. Create localbot.agent file with your corresponding Reddit credentials:
```
user_agent: "Reddit-Comments-Utils:<version> (by /u/<user>)"
client_id: "<client_id>"
client_secret: "<client_secret>"
username: "<username>"
password: "<password>"
```
## Run Instructions
### Event Listener (Requires RabbitMQ/PostGres setup)
Modify/run .bat files as needed with default settings
OR
Run manually:

Event listener
```
export GRAW_BOT_AGENT=localbot.agent && go run src/events/exec/main.go --subreddits <subreddits comma separated> --search <search text inside post body> --queue <queue to publish per amqp>
```

Consumer
```
go run src/consumer/exec/main.go --queue <queue to consume per amqp>
```

Can set environment variables as needed:
- GRAW_BOT_AGENT: location of your graw config
- AMQP_DEFAULT_BROKER: your broker url amqp://<user>:<pass>@localhost:5672
- AMQP_DEFAULT_EXCHANGE: amqp exchange
- POSTGRES_CONN: Postgres connection info

### Batch
To export multiple subreddit threads based on a regex "form style" search into *multiple* .csv files:
```
export GRAW_BOT_AGENT=localbot.agent && go run src/batch/exec/main.go csv --subreddit OMSA --threads a73ni1,8m2anv --headers "Application Date,Decision Date,Education,Test Scores,Experience,Recommendations,Comments"
```

To export multiple subreddit threads based on a regex "form style" search into a *single* .csv format:
(Notice the -m merge flag)
```
export GRAW_BOT_AGENT=localbot.agent && go run src/batch/exec/main.go csv --subreddit OMSA --threads a73ni1,8m2anv --headers "Application Date,Decision Date,Education,Test Scores,Experience,Recommendations,Comments" -m
```

To export multiple subreddit threads into a single .txt format:
```
export GRAW_BOT_AGENT=localbot.agent && go run src/batch/exec/main.go txt --subreddit OMSA --threads a73ni1,8m2anv
```

## Regex "Form Style" Search
To use the regex search, a Reddit thread must have replies that can be searched/filtered.
For example, the following "form" would work:
```
Application Date: 02/16/2019

Status: Won!

Comments: Hey!
```
Replace the <HEADER:> sections with the headers you need

## Tests

## Docker Compose

To use docker compose:
1. Create your localbot.agent file inside of src/events/exec/
2. Run:
```
docker-compose build
docker-compose up
```
The services for rabbitmq, postgres, consumer, and producer builds should initialize

**Obviously you should update your credentials using secrets and not use "admin:admin" for everything**
