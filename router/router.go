package router

import (
	"gateway/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	//初始化路由
	router := gin.Default()

	// 性能分析
	//pprof.Register(router)

	// 记录请求日志
	router.Use(middleware.LogMiddleware())
	// 替换 gin 的 recovery 处理方法
	router.Use(middleware.LogRecoveryMiddleware(false))
	// default request
	router.GET("/", func(context *gin.Context) {
		//panic("bbb")
		context.Writer.Write([]byte("hello world!!!"))
	})

	return router
}
