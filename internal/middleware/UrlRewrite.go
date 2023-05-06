package middleware

import (
	"github.com/gin-gonic/gin"
)

//url重写
func HttpUrlRewriteMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//proxyServer, ok := ctx.Get("service")
		//if !ok {
		//	zap.L().Info("未能从上下文中获取该服务详细信息")
		//	ctx.Abort()
		//	return
		//}

		//serviceDetail := proxyServer.((*loadBalance.ServerDetail))

		//TODO 匹配前 匹配后,匹配前 匹配后  多个重写规则
		//127.0.0.1:9090:/test/ab->127.0.0.1:9090:/test/ba

		ctx.Next()
	}
}

//stripuri功能,和urlrewrite做好冲突
//访问网关127.0.0.1:9090/test/ab->期望的下游(上游服务器)的地址(实际的服务地址)127.0.0.1:9900/ab,stripuri清空多余的前缀
func HttpStripUriMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//proxyServer, ok := ctx.Get("service")
		//if !ok {
		//	zap.L().Info("未能从上下文中获取该服务详细信息")
		//	ctx.Abort()
		//	return
		//}

		//serviceDetail := proxyServer.((*loadBalance.ServerDetail))
		//前缀接入方式

		//TODO 替换当前链接
		//ctx.Request.URL.Path = //strings.Replace(ctx.Request.URL.Path, serviceDetail.Http.Rule, "", 1) //替换一次
		//
		//ctx.Next()
	}
}
