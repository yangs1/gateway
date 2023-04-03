package middleware

import (
	"gateway/global"
	"gateway/initialize"
	"github.com/gin-gonic/gin"
	"net/http/httputil"
	"net/url"
)

// 根据http请求类型，选择负载均衡方式
func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lb := initialize.ServiceManagerHandler.GetLoadBalancer("http_proxy")

		if lb == nil {
			global.Logger.Error("获取负载均衡器失败")
			ctx.Abort()
			return
		}

		pryUrl, err := lb.Get(ctx.Request.RemoteAddr)

		if err != nil || pryUrl == "" {
			panic("get next addr fail")
		}

		target, err := url.Parse(pryUrl)

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(ctx.Writer, ctx.Request) //提供http服务(也是http服务暴露出去)
		ctx.Abort()
		return
	}
}
