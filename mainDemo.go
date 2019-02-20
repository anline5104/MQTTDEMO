package main

import (
	"demo/mqtt"
	"demo/tools"
	"encoding/json"
	"log"
	"os"
)

type SrcTopic struct {
	Src           string   `json:"src"`
	Topic         string   `json:"topic"`
}

type ServersDemo struct {
	Dest          string    `json:"dest"`
	Topics        []SrcTopic
}

//对应队列
var descToTopic []map[string](map[string]bool)   //唯一的desc 对应topic
var srcToTopic  []map[string](map[string]bool)   //唯一的src  对应topic

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

	//添加进对应desc队列
	for k,_ := range distinctDestMap{
		for i:=0;i< len(serversDemos);i++{
			destbroker := "tcp://" + serversDemos[i].Dest
			if k == destbroker {
				for j:=0;j< len(serversDemos[i].Topics);j++{
					descTopic :=  serversDemos[i].Topics[j].Topic

					descTo := make(map[string]map[string]bool)
					topic := make(map[string]bool)
					topic[descTopic] = true

					descTo[k] = topic
					descToTopic = append(descToTopic,descTo)
				}
			}
		}
	}

	//添加进对应的src队列
	for k,_ := range distinctSrcMap{
		for i:=0;i< len(serversDemos);i++{
				for j:=0;j< len(serversDemos[i].Topics);j++{
					srcbroker :=  "tcp://" + serversDemos[i].Topics[j].Src
					if k == srcbroker{
						srcTopic :=  serversDemos[i].Topics[j].Topic

						srcTo := make(map[string]map[string]bool)
						topic := make(map[string]bool)
						topic[srcTopic] = true
						srcTo[k] = topic
						srcToTopic = append(srcToTopic,srcTo)
					}
				}
			}
		}

	//数组去重处理后的队列
	descToTopicNoRep := tools.RemoveRepByLoop(descToTopic)  //map[tcp://inspect.huayuan-iot.com:map[sample-values/9H200A1700023/#:true]]
	srcToTopicNoRep := tools.RemoveRepByLoop(srcToTopic)


	//分类client
	destClientMap := make(map[string]mqtt.MQTTClient)
	srcClientMap  := make(map[string]mqtt.MQTTClient)


	//发布
	for k,_ := range distinctDestMap{
		destClientMap[k] = mqtt.DestInit(k)
		destcli := destClientMap[k]
		go destcli.Connect()

		for _,v := range descToTopicNoRep{
			for key,value := range v{
				if key == k {
					for keyTopic,_ := range value{
						go func(M mqtt.MQTTClient,s string) {
							<-destcli.Connected
							destcli.PublishSampleValues(keyTopic,"helloWorld")
						}(destcli,keyTopic)
					}
				}
			}
		}
	}

	//订阅
	for k,_ := range distinctSrcMap{
		srcClientMap[k] = mqtt.SrcInit(k)
		srccli  :=  srcClientMap[k]
		go srccli.Connect()

		for _,v := range srcToTopicNoRep{
			for key,value := range v{
				if key == k {
					for keyTopic,_ := range value{
						go func(M mqtt.MQTTClient,s string) {
							<-srccli.Connected
							srccli.Subscribe(keyTopic)
						}(srccli,keyTopic)
					}
				}
			}
		}
	}

}