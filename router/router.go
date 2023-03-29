package router

import (
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	//初始化路由
	router := gin.Default()

	router.GET("/", func(context *gin.Context) {
		context.Writer.Write([]byte("hello world!!!"))
	})

	return router
}

