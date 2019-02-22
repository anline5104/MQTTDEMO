package main

import (
	"demo/mqtt"
		"encoding/json"
	"log"
	"os"
		"strings"
	"demo/tools"
	"os/signal"
	"syscall"
)

type SrcTopic struct {
	Src           string   `json:"src"`
	Topic         string   `json:"topic"`
}

type ServersDemo struct {
	Dest          string    `json:"dest"`
	Topics        []SrcTopic
}

//对应map
var destToTopic map[string](map[string]bool)   //唯一的dest 对应topic
var srcToTopic  map[string](map[string]bool)   //唯一的src  对应topic


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
	destTopic := make(map[string]bool) //dest的topic map
	destToTopic = make(map[string](map[string]bool))

	for k,_ := range distinctDestMap{
		for i:=0;i< len(serversDemos);i++{
			destbroker := "tcp://" + serversDemos[i].Dest
			if k == destbroker {
				for j:=0;j< len(serversDemos[i].Topics);j++{
					//dest的topic map
					destTopicSingle := serversDemos[i].Topics[j].Topic
					destTopic[destTopicSingle] = true
				}
				destToTopic[k] = destTopic
			}
		}
	}

	//添加进对应的src队列
	srcToTopic = make(map[string](map[string]bool))
	for k,_ := range distinctSrcMap{
		for i:=0;i< len(serversDemos);i++{
			for j:=0;j< len(serversDemos[i].Topics);j++{
				srcbroker :=  "tcp://" + serversDemos[i].Topics[j].Src
				srcTopicSingele := serversDemos[i].Topics[j].Topic
				if k == srcbroker{
					srcTopic := make(map[string]bool)
					srcTopic[srcTopicSingele] = true
					srcToTopic[k] = srcTopic
					//srcToTopic = append(srcToTopic,srcTo)
				}
			}
		}
	}

	//数组去重处理后的队列
	//descToTopicNoRep := tools.RemoveRepByLoop(descToTopic)  //map[tcp://inspect.huayuan-iot.com:map[sample-values/9H200A1700023/#:true]]
	//srcToTopicNoRep := tools.RemoveRepByLoop(srcToTopic)


	//分类client
	destClientMap := make(map[string]mqtt.MQTTClient)
	srcClientMap  := make(map[string]mqtt.MQTTClient)


	//连接Dest
	for k,_ := range distinctDestMap{
		destcli := mqtt.DestInit(k)
		destClientMap[k] = destcli
		go destcli.Connect()
	}

	log.Println("SABnn")

	//连接src
	for k,_ := range distinctSrcMap{
		srcClientMap[k] = mqtt.SrcInit(k)
		srccli  :=  srcClientMap[k]
		go srccli.Connect()


		for key,value := range srcToTopic{
			if key == k {
				for keyTopic,_ := range value{
					go func( ) {
						<-srccli.Connected
						srccli.Client.Subscribe(keyTopic,0,nil)
					}()
				}
			}
		}

	}
   /* time.Sleep(time.Second * 10)*/

	//发布
RetryLoop:
	for {
		<-mqtt.MsgChan
		log.Println(mqtt.MsgQueue)
		for _, ms := range mqtt.MsgQueue {                 //遍历msg的队列
			mTopic := ms.Topic()
			mPayload := ms.Payload()
			for destUrl, destClient := range destClientMap {          //遍历destClientMap，其中destUrl为tcp://inspect.huayuan-iot.com
					for url, topicMap := range destToTopic {       //遍历destToTopicMap,url为tcp://inspect.huayuan-iot.com，topicMap为 map[sample-values/9H200A1700023/#:true]
						if (url == destUrl) {
							for topicKey, _ := range topicMap {
								if (strings.Contains(mTopic,tools.CutOutString(topicKey))) {
									go func(destClient *mqtt.MQTTClient,mTopic string,mPayload []byte) {
										<-destClient.Connected
										destClient.Client.Publish(mTopic, 0, true, mPayload)
									}(&destClient,mTopic,mPayload)
								}
							}
						}
					}

			}
		}

		exit := make(chan os.Signal,10) //初始化一个channel
		signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM) //notify方法用来监听收到的信号
		select {
		case <-exit:
			break RetryLoop
		}
	}

}