package test

import (
	"testing"
	"demo/tools"
	"log"
)

func  TestString(t *testing.T){
	log.Println(tools.CutOutString("sample-values/9H200A1700008/#"))

}