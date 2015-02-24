package server

import (
	"reflect"
	"testing"

	"golang.org/x/net/websocket"
)

type TestConn struct {
	in  interface{}
	out interface{}
	ws  *websocket.Conn
}

func (c *TestConn) ReadJSON(msg interface{}) error {
	in := c.in.(Message)
	rv := reflect.ValueOf(msg).Elem()
	rv.Set(reflect.ValueOf(&in).Elem())
	return nil
}

func (c *TestConn) WriteJSON(msg interface{}) error {
	c.out = msg
	return nil
}

func TestTestReadJSON(t *testing.T) {
	out := Message{}
	exp := Message{ID: "test"}
	conn := TestConn{in: exp}
	conn.ReadJSON(&out)
	if out.ID != exp.ID {
		t.Error("should set test value")
	}
}
