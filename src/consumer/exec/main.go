package main

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/dillonmabry/reddit-comments-util/src/datamanager"
	"github.com/dillonmabry/reddit-comments-util/src/distributed"
	"github.com/dillonmabry/reddit-comments-util/src/logging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var logger = logging.NewLogger()

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	var postMessage datamanager.PostMessage
	if err := json.Unmarshal(msg.Payload(), &postMessage); err != nil {
		logger.Error("Error converting payload to message struct")
	}
	exists, err := datamanager.IsPostExist(postMessage.URL)
	if err != nil {
		logger.Error("Error checking if post exists")
	}
	if exists != true {
		datamanager.SavePostMessage(&postMessage)
	}
}

func main() {
	datamanager.InitDB("host=localhost port=5432 user=admin password=admin dbname=reddit sslmode=disable")
	topic := flag.String("topic", "", "Specifies topic to consume and persist posts")
	flag.Parse()
	c := distributed.NewDistributed("tcp://192.168.1.220:1883", *topic, f)

	if token := c.Client.Subscribe(c.Topic, 0, nil); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}

	done := make(chan bool)
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	<-done
}
