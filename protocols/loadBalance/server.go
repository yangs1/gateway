package loadBalance

import "gateway/model/loadBalance"

type Server struct {
	LbHandler []*ServerDetail
}

type ServerDetail struct {
	BaseInfo    *loadBalance.ServerInfo
	LoadBalance *LoadBalanceNode
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
