package main

import (
	"gateway/initialize"
	"gateway/provider"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	initialize.InitLogger()

	go func() {
		provider.StartHttpServer()
	}()
	// 当前的goroutine等待信号量
	quit := make(chan os.Signal)
	// 监控信号：SIGINT, SIGTERM, SIGQUIT
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	// 这里会阻塞当前goroutine等待信号
	<-quit

	provider.StopHttpServer()

}
