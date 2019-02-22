package mqtt

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"time"
)

type MQTTClient struct {
	Client       MQTT.Client
	Name         string
	Connected    chan bool
	Hasconnected bool
}

//消息队列
var MsgQueue []MQTT.Message
var MsgChan = make(chan bool)

//init dest
func DestInit(destbroker string) MQTTClient{
	destopts := MQTT.NewClientOptions()
	destopts.AddBroker(destbroker)
	destopts.SetClientID("ShareDestinationClientID" + time.Now().String())
	destopts.SetCleanSession(false)
	destopts.SetStore(MQTT.NewFileStore("memory"))
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
	destcli.Connected = make(chan bool)
	destcli.Hasconnected = false

	log.Println(destcli)

	return destcli
}

//init src
func SrcInit(srcbroker string) MQTTClient{
	srcopts := MQTT.NewClientOptions()
	srcopts.AddBroker(srcbroker)
	srcopts.SetClientID("ShareSourceClientID"+ time.Now().String())
	srcopts.SetCleanSession(false)
	srcopts.SetStore(MQTT.NewFileStore("memory"))
	srcopts.SetAutoReconnect(true)
	srcopts.SetOnConnectHandler(func(client MQTT.Client) {
		log.Printf("source service [%v] has connected\n", srcbroker)
	})
	srcopts.SetConnectionLostHandler(func(client MQTT.Client, err error) {
		log.Printf("source service [%v] disconnected\n", srcbroker)
	})

	srcopts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		MsgQueue = append(MsgQueue,msg)
		MsgChan <- true
	})

	// share message


	/*srcopts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		// public message
		if destcli.Hasconnected == true {
			// log.Printf("public message: %v, %s", msg.Topic(), msg.Payload())

			token := destcli.Client.Publish(msg.Topic(), 0, true, msg.Payload())
			token.Wait()
		}
	})*/

	srcsrv := MQTT.NewClient(srcopts)
	var srccli MQTTClient
	srccli.Client = srcsrv
	srccli.Name = srcbroker
	srccli.Connected = make(chan bool)
	srccli.Hasconnected = false

	log.Println(srccli)
	return srccli
}


func (c *MQTTClient) Connect() error {
	retry := time.NewTicker(5 * time.Second)
RetryLoop:
	for {
		select {
		case <-retry.C:
			log.Println(c)
			token := c.Client.Connect()
			log.Println(token)
			if token.Wait() && token.Error() != nil {
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
	c.Connected <- true

	return nil
}


