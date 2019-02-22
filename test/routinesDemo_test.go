package test

import (
	"testing"
)

func  TestRoutines(t *testing.T){
	defer println("执行完毕")

	for i:=0;i<3;i++{
	go func() {
		for i:=100;i>0;i--{
			//time.Sleep(1*time.Second)
			println("开始打印%s",i)
		}
	}()
	}
    println("结束")
}
