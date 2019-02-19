package main

import (
	"flag"
	"log"
	"os"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient mqtt client
type MQTTClient struct {
	Client       MQTT.Client
	Name         string
	connected    chan bool
	hasconnected bool
}


func main() {
	srcPtr := flag.String("src", "cloud.huayuan-iot.com:1883", "source service")
	destPtr := flag.String("dest", "inspect.huayuan-iot.com:1883", "destination service") //
	topicPtr := flag.String("topic", "sample-values/9H200A1700008/#", "shared topic")

	flag.Parse()

	src := *srcPtr
	dest := *destPtr
	topic := *topicPtr

	srcbroker := "tcp://" + src
	destbroker := "tcp://" + dest

	// destination server init
	destopts := MQTT.NewClientOptions()
	destopts.AddBroker(destbroker)
	destopts.SetClientID("ShareDestinationClientID" + time.Now().String())
	destopts.SetCleanSession(false)
	destopts.SetStore(MQTT.NewFileStore(":memory:"))
	destopts.SetAutoReconnect(true)
	destopts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Printf("destination service [%v] has connected\n", destbroker)
	})
	destopts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		log.Printf("destination service [%v] disconnected\n", destbroker)
	})

	destsrv := MQTT.NewClient(destopts)

	var destcli MQTTClient
	destcli.Client = destsrv
	destcli.Name = destbroker
	destcli.connected = make(chan bool)
	destcli.hasconnected = false

	// source server init
	srcopts := MQTT.NewClientOptions()
	srcopts.AddBroker(srcbroker)
	srcopts.SetClientID("ShareSourceClientID")
	srcopts.SetCleanSession(false)
	srcopts.SetStore(MQTT.NewFileStore(":memory:"))
	srcopts.SetAutoReconnect(true)
	srcopts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Printf("source service [%v] has connected\n", srcbroker)
	})
	srcopts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		log.Printf("source service [%v] disconnected\n", srcbroker)
	})


	// share message
	srcopts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		// public message
		if destcli.hasconnected == true {
			// log.Printf("public message: %v, %s", msg.Topic(), msg.Payload())

			token := destcli.Client.Publish(msg.Topic(), 0, true, msg.Payload())
			token.Wait()
		}
	})

	srcsrv := MQTT.NewClient(srcopts)

	var srccli MQTTClient
	srccli.Client = srcsrv
	srccli.Name = srcbroker
	srccli.connected = make(chan bool)
	srccli.hasconnected = false

	// connect servers
	go destcli.connect()
	go srccli.connect()

	// source connected
	<-srccli.connected

	if token := srccli.Client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Printf("subscribe topic [%v] fail, exit!\n", topic)
		os.Exit(1)
	}

	<-destcli.connected
	destcli.hasconnected = true

	for {
		time.Sleep(time.Second)
	}
}

// connect server
func (c MQTTClient) connect() error {
	retry := time.NewTicker(5 * time.Second)
RetryLoop:
	for {
		select {
		case <-retry.C:
			if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
				// handle error
				log.Printf("connect mqtt server [%v] fail\n", c.Name)
			} else {
				// success
				log.Printf("connect mqtt server [%v] success\n", c.Name)
				retry.Stop()
				break RetryLoop
			}
		}
	}

	c.connected <- true

	return nil
}
