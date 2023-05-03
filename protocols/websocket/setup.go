package websocket

import "gateway/protocols"

const (
	Protocol = "websocket"
)

func init() {
	protocols.ServerIoc.Register(Protocol, setup)
}

func setup() protocols.Server {
	server := InitServer()

	return server
}
