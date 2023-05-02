package main

import (
	"fmt"
	"gateway/config"
	"gateway/initialize"
	"gateway/model"
	"gateway/model/loadBalance"
	"testing"
)

func TestDb(t *testing.T) {
	initialize.Viper(&config.GVA_CONFIG, "../config/app.debug.yaml")

	fmt.Println(config.GVA_CONFIG)

	config.GVA_LOG = initialize.Zap()
	config.GVA_DB = initialize.GormMysql()

	serviceInfo := &loadBalance.ServerInfo{}

	serverRes, serverNum, _ := serviceInfo.PageList(config.GVA_DB, model.PageInput{PageSize: 100, PageNum: 1})

	t.Log(serverRes, serverNum)

	for _, resModel := range serverRes {
		serviceHttp := &loadBalance.ServerHttp{}
		rrr, _ := serviceHttp.PageList(config.GVA_DB, resModel)

		t.Log(rrr)

	}

}
