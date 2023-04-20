package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter() *gin.Engine {
	//初始化路由
	router := gin.Default()

	// 性能分析
	//pprof.Register(router)

	/*router.Use(
		middleware.LogMiddleware(),              // 记录请求日志
		middleware.LogRecoveryMiddleware(false), // 替换 gin 的 recovery 处理方法
		middleware.HttpReverseProxyMiddleware(), //代理到实际的服务工作负载
	)*/

	// default request
	router.GET("/", func(context *gin.Context) {
		//panic("bbb")
		context.Writer.Write([]byte("hello world!!!"))
	})

	return router
}
