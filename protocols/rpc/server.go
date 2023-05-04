package rpc

import (
	"go.uber.org/zap"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
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

func InitServer() {
	server := &Server{}

	if err := rpc.Register(server); err != nil {
		log.Panic("Gateway api server run failed, error: "+err.Error(), zap.Error(err))
	}

	listener, err := net.Listen("tcp", "127.0.0.1:9099")
	if err != nil {
		log.Panic("Gateway api server run failed, error: "+err.Error(), zap.Error(err))
	}

	defer func() {
		if err := listener.Close(); err != nil {
			log.Panic("Gateway api server run failed, error: "+err.Error(), zap.Error(err))
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go jsonrpc.ServeConn(conn)
	}
}
