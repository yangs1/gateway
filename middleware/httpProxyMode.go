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

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(ctx.Writer, ctx.Request) //提供http服务(也是http服务暴露出去)
		ctx.Abort()
		return
	}
}

//新建一个实际的http的代理服务器(处理http https请求),gin的middleware调用
/*func NewLoadBalanceReverseProxy(c *gin.Context, lb loadBalance.LoadBalance, trans *http.Transport) *httputil.ReverseProxy {
	//请求协调者
	//请求协调者
	director := func(req *http.Request) {
		nextAddr, err := lb.Get(req.URL.String())
		if err != nil || nextAddr == "" {
			panic("get next addr fail")
		}
		target, err := url.Parse(nextAddr)
		if err != nil {
			panic(err)
		}
		targetQuery := target.RawQuery
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "user-agent")
		}
	}

	//更改内容
	modifyFunc := func(resp *http.Response) error {
		if strings.Contains(resp.Header.Get("Connection"), "Upgrade") { //判断Upgrade是否在connection中,协议升级,ws
			return nil
		}
		return nil
	}

	//错误回调 ：关闭real_server时测试，错误回调
	//范围：transport.RoundTrip发生的错误、以及ModifyResponse发生的错误
	errFunc := func(w http.ResponseWriter, r *http.Request, err error) {
		//middleware.ResponseError(c, 999, err)
		global.Logger.Info("调用错误回调函数")
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyFunc, ErrorHandler: errFunc}
}
*/
