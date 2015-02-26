package server

import "golang.org/x/net/websocket"

type WSConn struct {
	ws *websocket.Conn
}

func (c *WSConn) ReadMessage() (Message, error) {
	msg := Message{}
	err := websocket.JSON.Receive(c.ws, &msg)
	return msg, err
}

func (c *WSConn) WriteMessage(msg interface{}) error {
	return websocket.JSON.Send(c.ws, msg)
}
