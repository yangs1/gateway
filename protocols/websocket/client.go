package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 10 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//Client is websocket client
type Client struct {
	clientId string
	conn     *websocket.Conn
	send     chan []byte

	server *Server
}

//Close client connection
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		zap.L().Error("Close err: "+err.Error(), zap.Error(err))
	}
	return
}

func (c *Client) Read() {
	defer func() {
		c.server.unregister <- c
	}()

	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		zap.L().Error("SetReadDeadline err: "+err.Error(), zap.Error(err))
		return
	}
	c.conn.SetPongHandler(func(string) error {
		if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			zap.L().Error("SetReadDeadline err: "+err.Error(), zap.Error(err))
		}
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		//_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				zap.L().Error("IsUnexpectedCloseError err: "+err.Error(), zap.Error(err))
			}
			break
		}
		fmt.Println("msg :" + string(message))
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()

		c.server.unregister <- c
		//if err := c.conn.Close(); err != nil {
		//	zap.L().Error("conn Close err: "+err.Error(), zap.Error(err))
		//}
	}()

	for {
		select {
		case message, ok := <-c.send:
			// 设置超时时间
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				zap.L().Error("SetWriteDeadline err: "+err.Error(), zap.Error(err))
			}
			if !ok {
				// The hub closed the channel.
				if err := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
					zap.L().Error("WriteMessage err: "+err.Error(), zap.Error(err))
				}
				return
			}
			//err := c.conn.WriteMessage(websocket.TextMessage, message)
			//if err != nil {
			//	zap.L().Error("WriteMessage err: "+err.Error(), zap.Error(err))
			//	return
			//}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				zap.L().Error("NextWriter err: "+err.Error(), zap.Error(err))
				return
			}
			// 写失败不结束连接
			if _, err := w.Write(message); err != nil {
				zap.L().Error("Write err: "+err.Error(), zap.Error(err))
			}

			if err := w.Close(); err != nil {
				zap.L().Error("Close err: "+err.Error(), zap.Error(err))
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				zap.L().Error("SetWriteDeadline err: "+err.Error(), zap.Error(err))
				return
			}
			// ping 失败直接断开
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				zap.L().Error("WriteMessage err: "+err.Error(), zap.Error(err))
				return
			}
		}
	}
}
