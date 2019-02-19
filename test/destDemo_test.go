package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)


func TestMarshalJson(t *testing.T){
	type SrcTopic struct {
		Src           string   `json:"src"`
		Topic         string   `json:"topic"`
	}

	type ServersDemo struct {
		Dest          string    `json:"dest"`
		Topics        []SrcTopic
	}


	filePtr,err := os.Open("../destDemo.json")
	if err != nil{
		fmt.Println("Open file failed")
		return
	}
	defer filePtr.Close()


	var ServersDemos []ServersDemo

	//Decoder json
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&ServersDemos)
    if err != nil{
    	fmt.Println("Decoder failed")
    	os.Exit(0)
	}else{
		fmt.Println(ServersDemos)
	}

	//ergodic dest
	for _,Server := range ServersDemos{
		fmt.Println(Server.Dest)
		for _,topic := range Server.Topics{
			fmt.Print(" "+topic.Src)
		}
		fmt.Println()
	}
}
