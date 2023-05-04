package rpc

import (
	"fmt"
	"gateway/protocols"
	"gateway/protocols/websocket"
)

//Message is api message
type Message struct {
	Connections []string `json:"connections"` //消息接受者
	//Token string `json:"token"` //作为消息发送鉴权
	Msg string `json:"msg"` //为一个json，里边包含 type 消息类型
}

//Response is api response
type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func (s *Server) success(data interface{}) Response {
	return Response{
		Code:  200,
		Data:  data,
		Error: "",
	}
}

func (s *Server) failure(code int, err string) Response {
	return Response{
		Code:  code,
		Data:  "",
		Error: err,
	}
}

//SendToConnections send message to connections
func (s *Server) SendToConnections(message *Message, reply *Response) error {

	if len(message.Connections) == 0 {
		*reply = s.failure(400, "empty connections")
		return nil
	}

	if message.Msg == "" {
		*reply = s.failure(400, "empty msg")
		return nil
	}

	websocketServer, _ := protocols.ServerIoc.Load(websocket.Protocol)

	for _, clientId := range message.Connections {
		fmt.Println("clientId :" + clientId)
		if err := websocketServer.(*websocket.Server).SendToConnection(clientId, message.Msg); err != nil {
			*reply = s.failure(400, err.Error())
			return nil
		}
	}

	*reply = s.success(nil)
	return nil
}
