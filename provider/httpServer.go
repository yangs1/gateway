package provider

import (
	"context"
	"gateway/global"
	"gateway/pkg/util"
	"gateway/router"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

var HttpServerHandler *http.Server

func StartHttpServer() error {
	vcfg, err := util.ReadConfigFile("app")
	if err != nil {
		global.Logger.Error("general配置文件读取失败", zap.Error(err))
		return err
	}

	r := router.Router()
	HttpServerHandler = &http.Server{
		Handler: r, //gin.Engine
		Addr:    vcfg.GetString("http.address"),
		//time.Duration单位纳秒
		ReadTimeout:    time.Duration(vcfg.GetInt("http.read_timeout")) * time.Second, //*second
		WriteTimeout:   time.Duration(vcfg.GetInt("http.write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << vcfg.GetInt("http.max_header_bytes"),
	}

	if err := HttpServerHandler.ListenAndServe(); err != nil {
		global.Logger.Error("http server start error.", zap.Error(err))
	}

	return nil
}

func StopHttpServer() {
	// 调用Server.Shutdown graceful结束
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := HttpServerHandler.Shutdown(timeoutCtx); err != nil {
		log.Println("Server Shutdown:", err)
	}
}
