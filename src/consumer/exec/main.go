package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/dillonmabry/reddit-comments-util/src/datamanager"
	"github.com/dillonmabry/reddit-comments-util/src/distributed"
	"github.com/dillonmabry/reddit-comments-util/src/logging"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
	var postMessage datamanager.PostMessage
	if err := json.Unmarshal(msg.Payload(), &postMessage); err != nil {
		log.Println("Error unmarshalling payload to message struct")
	}
	datamanager.SavePostMessage(postMessage)
}

func main() {
	logger := logging.NewLogger()
	datamanager.InitDB("host=localhost port=5432 user=admin password=admin dbname=reddit sslmode=disable")
	c := distributed.NewDistributed("tcp://192.168.1.220:1883", "topic/test", f)

	if token := c.Client.Subscribe(c.Topic, 0, nil); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}

	if token := c.Client.Unsubscribe("some_topic"); token.Wait() && token.Error() != nil {
		logger.Fatal(token.Error())
	}
	//TODO: Cleanup forever loop
	done := make(chan bool)
	go func() {
		for {
			time.Sleep(time.Second)
		}
	}()
	<-done
}
