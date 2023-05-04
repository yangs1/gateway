package main

import (
	"fmt"
	"gateway/config"
	"gateway/initialize"
	"gateway/protocols/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

func TestRpc(t *testing.T) {
	initialize.Viper(&config.GVA_CONFIG, "../config/app.debug.yaml")

	fmt.Println(config.GVA_CONFIG)

	config.GVA_LOG = initialize.Zap()
	config.GVA_DB = initialize.GormMysql()

	rpca, err := jsonrpc.Dial("tcp", "127.0.0.1:9099")
	if err != nil {
		t.Error(err)
	}

	//response 为 json 字符串
	var response rpc.Response
	//调用远程方法
	//注意第三个参数是指针类型

	//发送消息
	err2 := rpca.Call("Server.SendToConnections", &rpc.Message{
		Connections: []string{"1"},
		Msg:         "hello rpc",
	}, &response)

	if err2 != nil {
		t.Error(err2)
	}

	fmt.Println(response)

}
