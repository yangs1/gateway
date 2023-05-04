package websocket

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

type Server struct {
	upgrader websocket.Upgrader

	clients sync.Map //map[string]*Client

	register   chan *Client
	unregister chan *Client
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

// Run start http server.
func InitServer() *Server {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	// 初始化 websocket 配置
	s := &Server{
		clients: sync.Map{},

		upgrader:   upgrader,
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	// 监听用户注册map
	go s.handleClients()

	return s
}

// 用户注册服务
func (s *Server) handleClients() {
	for {
		select {
		case client := <-s.register:
			s.clients.Store(client.clientId, client)
			//s.clients[client.clientId] = client

			//notify router api server
			//todo
			//if _, err := s.rpc.Call("Online", client.clientId, s.config.Node); err != nil {
			//	log.Error("Client online notification failed: %s" + err.Error(), zap.Error(err))
			//}
		case client := <-s.unregister:
			if _, ok := s.clients.LoadAndDelete(client.clientId); ok {
				//delete(s.clients, client.clientId)
				close(client.send)
				client.Close()

				//notify router api server
				//todo
				//if _, err := s.rpc.Call("Offline", client.connId, s.config.Node); err != nil {
				//	log.Error("Client online notification failed: %s", err.Error())
				//}
			}
		}
	}
}

func (s *Server) ServeWs(context *gin.Context) {
	//panic("bbb")
	conn, err := s.upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		zap.L().Error("Websocket upgrade failed: %s"+err.Error(), zap.Error(err))
		return
	}

	client := &Client{
		conn:     conn,
		send:     make(chan []byte, 1024),
		clientId: "1",
	}

	s.register <- client

	go client.Write()
	go client.Read()
}

func (s *Server) SendToConnection(connId string, msg string) error {
	if client, ok := s.clients.Load(connId); ok {
		select {
		case client.(*Client).send <- []byte(msg):
			//zap.L().Info("SendToConnection " + connId + ": " + msg)
			return nil
		default:
			client.(*Client).Close()
			return errors.New("send message failed to " + connId)
		}
	}

	return errors.New("send message failed, connection: " + connId + " not found")
}
