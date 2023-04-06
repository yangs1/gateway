package middleware

import (
	"gateway/global"
	"gateway/initialize"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// 根据http请求类型，选择负载均衡方式
func HttpReverseProxyMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		lb := initialize.ServiceManagerHandler.GetLoadBalancer()

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

		proxy := httputil.NewSingleHostReverseProxy(target) //NewLoadBalanceReverseProxy(target)
		proxy.ServeHTTP(ctx.Writer, ctx.Request)            //提供http服务(也是http服务暴露出去)
		ctx.Abort()
		return
	}
}

//处理一些自定义配置的 ReverseProxy
func NewLoadBalanceReverseProxy(target *url.URL) *httputil.ReverseProxy {

	proxyRev := httputil.NewSingleHostReverseProxy(target)

	proxyRev.Transport = initialize.ServiceManagerHandler.GetTrans()

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		//middleware.ResponseError(c, 999, err)
		global.Logger.Info("调用错误回调函数")
	}
	proxyRev.ErrorHandler = errFunc

	return proxyRev
}

