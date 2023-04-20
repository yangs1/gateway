package http

import (
	protocols "gateway/protocols"
)

const (
	Protocol = "http"
)

func init() {
	protocols.ServerIoc.Register(Protocol, setup)
}

func setup() protocols.Server {
	server := &Server{}

	return server
}
