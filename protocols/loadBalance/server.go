package loadBalance

import (
	"gateway/config"
	"gateway/model"
	"gateway/model/loadBalance"
)

type Server struct {
	LbHandler []*ServerDetail
}

type ServerDetail struct {
	BaseInfo    *loadBalance.ServerInfo
	LoadBalance *LoadBalanceHandler
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
			LoadBalance: lb,
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
