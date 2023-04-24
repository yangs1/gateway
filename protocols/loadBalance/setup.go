package loadBalance

import "gateway/protocols"

const (
	Protocol = "loadBalance"
)

func init() {
	protocols.ServerIoc.Register(Protocol, setup)
}

func setup() protocols.Server {
	server := &Server{}

	return server
}
