package middleware

import (
	"gateway/protocols"
	"gateway/protocols/loadBalance"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

//处理http请求接入方式,前缀(路径)和域名
func HttpProxyModeMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		manager, err := protocols.ServerIoc.Load(loadBalance.Protocol)
		if err != nil {
			zap.L().Error("make load balance server error. ", zap.Error(err))
			ctx.Abort() //取消,不在向下传递
			return
		}

		var lbServer *loadBalance.ServerDetail

		for _, handler := range manager.(*loadBalance.Server).LbHandler {
			if handler.BaseInfo.ServerType != loadBalance.TYPE_HTTP_SERVER {
				continue
			}

			if handler.BaseInfo.AccessType == loadBalance.TYPE_ACCESS_DOMAIN {
				host := ctx.Request.Host
				host = host[0:strings.Index(host, ":")] //会连端口一块获取,返回str(:)在host中第一次出现的位置(字符串的索引包括空格)(如果找不到则返回-1；如果str为空，则返回0)

				if handler.BaseInfo.AccessRule == host { //匹配成功,找到要请求的服务,返回该服务的servicedetail
					lbServer = handler
					break
				}
			}

			if handler.BaseInfo.AccessType == loadBalance.TYPE_ACCESS_PRE_URL {
				//测试path的前缀是不是以serviceItem.HTTPRule.Rule开头
				path := ctx.Request.URL.Path
				if strings.HasPrefix(path, handler.BaseInfo.AccessRule) {
					lbServer = handler
					break
				}
			}

			if handler.BaseInfo.AccessType == loadBalance.TYPE_ACCESS_ALL {
				lbServer = handler
				break
			}
		}

		lbServer.HealthCheck()
		ctx.Set("proxy_server", lbServer) //传递ServerDetail
		ctx.Next()                        //传递给子context
	}
}
