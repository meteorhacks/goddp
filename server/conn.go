package server

import "golang.org/x/net/websocket"

type Conn struct {
	ws *websocket.Conn
}

func (c *Conn) ReadJSON(msg interface{}) error {
	return websocket.JSON.Receive(c.ws, msg)
}

func (c *Conn) WriteJSON(msg interface{}) error {
	return websocket.JSON.Send(c.ws, msg)
}
