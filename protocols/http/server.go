package http

import (
	"context"
	"gateway/config"
	"gateway/config/autoload"
	"gateway/internal/router"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Server struct {
	HttpServerHandler *http.Server
}

func (s *Server) IsBind() bool {
	return false
}

func (s *Server) Protocol() string {
	return Protocol
}

func (s *Server) Config() interface{} {
	return config.GVA_CONFIG.HTTP
}

// Run start http server.
func (s *Server) Run() {
	httpCfg := s.Config().(autoload.Http)
	r := router.RegisterRouter()

	s.HttpServerHandler = &http.Server{
		Handler: r, //gin.Engine
		Addr:    httpCfg.Address,
		//time.Duration单位纳秒
		ReadTimeout:    time.Duration(httpCfg.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(httpCfg.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << httpCfg.MaxHeaderBytes,
	}

	// TODO ListenAndServeTLS()
	if err := s.HttpServerHandler.ListenAndServe(); err != nil {
		zap.L().Error("http server start error.", zap.Error(err))
	}
}

// Stop close http server
func (s *Server) Stop() {
	//超时即关闭
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.HttpServerHandler.Shutdown(ctx); err != nil {
		zap.L().Info("https 代理服务器 HttpsProxyServer Stop error:%v\n", zap.Error(err))
	}
	zap.L().Info("https 代理服务器 HttpsProxyServer Stop\n")
}
