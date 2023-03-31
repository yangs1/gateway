package provider

import "net/http"

type GatewayEngine struct {
	HttpServerHandler *http.Server
}

func NewGatewayEngine() *GatewayEngine {
	return &GatewayEngine{}
}
