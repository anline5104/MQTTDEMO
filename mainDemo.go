package main

import (
	"demo/mqtt"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type SrcTopic struct {
	Src           string   `json:"src"`
	Topic         string   `json:"topic"`
}

type ServersDemo struct {
	Dest          string    `json:"dest"`
	Topics        []SrcTopic
}





func main() {
	filePtr, err := os.Open("destDemo.json")
	if err != nil {
		log.Fatalf("open `%s` file error: %s", "destDemo.json", err)
		return
	}
	defer filePtr.Close()

	var serversDemos []ServersDemo

	//Decoder json
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&serversDemos)
	if err != nil {
		log.Fatalf("Decoder failed: `%s`, error: %s", "destDemo.json", err)
		os.Exit(0)
	}

	//去除相同url
	distinctDestMap := make(map[string]bool)
	distinctSrcMap  := make(map[string]bool)

	for i:=0;i< len(serversDemos);i++{
		destbroker := "tcp://" + serversDemos[i].Dest
		distinctDestMap[destbroker] = true

		for j:=0;j< len(serversDemos[i].Topics);j++{
			srcbroker :=  "tcp://" + serversDemos[i].Topics[j].Src
			distinctSrcMap[srcbroker] = true
		}
	}

	//分类client
	destClientMap := make(map[string]mqtt.MQTTClient)
	srcClientMap  := make(map[string]mqtt.MQTTClient)


	for k,_ := range distinctDestMap{
		destClientMap[k] = mqtt.DestInit(k)
	}

	for k,_ := range distinctSrcMap{
		srcClientMap[k] = mqtt.SrcInit(k)
	}



	//遍历
	for _,server := range serversDemos{
		fmt.Println(server.Dest)
		destbroker :=  "tcp://" + server.Dest
		destcli := mqtt.DestInit(destbroker)
		go destcli.Connect()
	/*	<-destcli.Connected
		destcli.Hasconnected = true*/

		for _,topic := range server.Topics{
			srcbroker :=  "tcp://" + topic.Src
			srccli  :=  mqtt.SrcInit(srcbroker,destcli)

			go srccli.Connect()

			<-srccli.Connected

			if token := srccli.Client.Subscribe(topic.Topic, 0, nil); token.Wait() && token.Error() != nil {
				log.Printf("subscribe topic [%v] fail, exit!\n", topic)
				os.Exit(1)
			}

			for {
				time.Sleep(time.Second)
			}
		}

	}

}