package main

import (
	"fmt"
	"gateway/config"
	"gateway/initialize"
	"gateway/model"
	lbmodel "gateway/model/loadBalance"
	"gateway/protocols"
	"gateway/protocols/loadBalance"
	"testing"
)

func TestDb(t *testing.T) {
	initialize.Viper(&config.GVA_CONFIG, "../config/app.debug.yaml")

	fmt.Println(config.GVA_CONFIG)

	config.GVA_LOG = initialize.Zap()
	config.GVA_DB = initialize.GormMysql()

	serviceInfo := &lbmodel.ServerInfo{}

	serverRes, serverNum, _ := serviceInfo.PageList(config.GVA_DB, model.PageInput{PageSize: 100, PageNum: 1})

	t.Log(serverRes, serverNum)

	for _, resModel := range serverRes {
		serviceHttp := &lbmodel.ServerHttp{}
		rrr, _ := serviceHttp.PageList(config.GVA_DB, resModel)

		t.Log(rrr)

	}

	// protocols 测试
	aaaa, _ := protocols.ServerIoc.Load(loadBalance.Protocol)
	t.Log(aaaa.(*loadBalance.Server).LbHandler[0].LoadBalance)
	t.Log(aaaa.(*loadBalance.Server).LbHandler[0].BaseInfo)

}
