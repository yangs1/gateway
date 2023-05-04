package main

import (
	"fmt"
	"gateway/config"
	"gateway/initialize"
	"gateway/protocols"
	"gateway/protocols/http"
	"gateway/protocols/rpc"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.Viper(&config.GVA_CONFIG)

	fmt.Println(config.GVA_CONFIG)

	config.GVA_LOG = initialize.Zap()
	config.GVA_DB = initialize.GormMysql()
	//config.GVA_LOG.Error("test", zap.String("sad", "dasd"))
	//config.GVA_LOG.Info("info_test", zap.String("sad", "dasd"))
	//config.GVA_DB.Table("relations").Where("id = ? and owner_id=?", "23", "1").Rows()

	if engine, err := protocols.ServerIoc.Load(http.Protocol); err == nil {
		go func() {
			engine.(*http.Server).Run()
		}()

		// 启动rpc 服务
		go rpc.InitServer()

		// 当前的goroutine等待信号量
		quit := make(chan os.Signal)
		// 监控信号：SIGINT, SIGTERM, SIGQUIT
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 这里会阻塞当前goroutine等待信号
		<-quit

		engine.(*http.Server).Stop()
	} else {
		log.Fatal("http server start error:", err)
	}

	log.Println("http server close.")
}
