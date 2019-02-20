package test

import (
	"demo/tools"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

func  TestSlice(t *testing.T){
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

	filePtr, err := os.Open("../destDemo.json")
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

	fmt.Println("+++++ 显示去重后的URL+++++")
	for k,_ :=range distinctDestMap{
		fmt.Println(k)
	}
	fmt.Println("---------------------")

	for k,_ :=range distinctSrcMap{
		fmt.Println(k)
	}

	fmt.Println("+++++ 显示去重后的URL+++++")

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




	for _,v :=range descToTopic{
		fmt.Println(v)
	}

	fmt.Println("------------------")

	//未去重
	for _,v :=range srcToTopic{
		fmt.Println(v)
	}

	fmt.Println("++++ 数组去重 ++++")

	for _,v :=range tools.RemoveRepByLoop(descToTopic){
		fmt.Println(v)
	}

	fmt.Println("------------------")
	for _,v :=range tools.RemoveRepByLoop(srcToTopic){
		fmt.Println(v)
	}


}