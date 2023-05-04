package middleware

import (
	"fmt"
	"gateway/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

//jwt验证,token认证
func HttpJwtAuthTokenMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//proxyServer, ok := ctx.Get("proxy")
		//if !ok {
		//	zap.L().Error("未能从上下文中获取服务详细信息")
		//	ctx.Abort()
		//	return
		//}
		//if proxyServer.(*loadBalance.ServerDetail).BaseInfo.OpenAuth == 0 {
		//	ctx.Next()
		//	return
		//}

		authToken := ctx.Request.Header.Get("Authorization")

		//token := strings.ReplaceAll(authToken, "Basic ", "") //返回s的副本,new替换old->token
		auth_list := strings.Split(authToken, " ")
		if len(auth_list) == 2 {
			// && auth_list[0] == "Basic"
			// Basic 解码方式
			//res, err := base64.StdEncoding.DecodeString(auth_list[1])
			//log.Printf(string(res))
			//if err == nil && string(res) == "sx:123" {
			//	ctx.Next()
			//	return
			//}

			claims, err := utils.JwtDecode(auth_list[1]) //*jwt.StandardClaims,标准token字段(访问群体,过期时间,id,发行时间,发行机构)
			if err != nil {
				zap.L().Error("解码token失败")
				ctx.Abort()
				return
			}

			if claims.Id != "" {
				fmt.Println("get claims id:", claims.Id)
				ctx.Set("user_id", claims.Id) //设置租户信息进context
			}
		} else {
			zap.L().Error("鉴权失败")
			ctx.Abort()
			return
		}

	}
}
