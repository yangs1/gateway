package loadBalance

import (
	"fmt"
	"gateway/config"
	"gateway/model"
	"gateway/model/loadBalance"
	"sync"
)

//接入方式
const (
	TYPE_HTTP_SERVER = iota //http服务类型
	TYPE_TCP_SERVER         //tcp 服务类型
)

//访问方式
const (
	TYPE_ACCESS_ALL = iota
	TYPE_ACCESS_DOMAIN
	TYPE_ACCESS_PRE_URL
)

type Server struct {
	LbHandler []*ServerDetail
}

type ServerDetail struct {
	BaseInfo    *loadBalance.ServerInfo
	loadBalance *LoadBalanceHandler
	init        sync.Once
}

// 健康检查
func (s *ServerDetail) HealthCheck() {
	s.init.Do(func() {
		fmt.Println("gogogogogo")
		s.loadBalance.Watch()
	})
}

func NewServer() *Server {

	server := &Server{}

	serviceInfo := &loadBalance.ServerInfo{}

	curDb := config.GVA_DB

	serverLists, _, _ := serviceInfo.PageList(curDb, model.PageInput{PageSize: 100, PageNum: 1})

	for _, infoModel := range serverLists {
		lb := NewLoadBalance(infoModel.RoundType)
		serverHttp := &loadBalance.ServerHttp{}
		httpLists, _ := serverHttp.PageList(curDb, infoModel)

		for _, httpModel := range httpLists {
			lb.AddNode(httpModel)
		}

		server.LbHandler = append(server.LbHandler, &ServerDetail{
			BaseInfo:    &infoModel,
			loadBalance: lb,
		})
	}

	return server
}

func (s *Server) IsBind() bool {
	return true
}

func (s *Server) Protocol() string {
	return Protocol
}

func (s *Server) Config() interface{} {
	return nil
}
