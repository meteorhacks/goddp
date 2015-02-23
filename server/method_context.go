package server

import (
	"errors"

	"golang.org/x/net/websocket"
)

type MethodContext struct {
	ID      string
	Args    []interface{}
	Conn    *websocket.Conn
	Done    bool
	Updated bool
}

func NewMethodContext(m Message, ws *websocket.Conn) MethodContext {
	ctx := MethodContext{}
	ctx.ID = m.ID
	ctx.Args = m.Params
	ctx.Conn = ws
	return ctx
}

func (ctx *MethodContext) SendResult(r interface{}) error {
	if ctx.Done {
		err := errors.New("already sent results for method")
		return err
	}

	ctx.Done = true
	return websocket.JSON.Send(ctx.Conn, map[string]interface{}{
		"msg":    "result",
		"id":     ctx.ID,
		"result": r,
	})
}

func (ctx *MethodContext) SendError(e string) error {
	if ctx.Done {
		err := errors.New("already sent results for method")
		return err
	}

	ctx.Done = true
	return websocket.JSON.Send(ctx.Conn, map[string]interface{}{
		"msg": "result",
		"id":  ctx.ID,
		"error": map[string]string{
			"error": e,
		},
	})
}

func (ctx *MethodContext) SendUpdated() error {
	if ctx.Updated {
		err := errors.New("already sent updated for method")
		return err
	}

	ctx.Updated = true
	return websocket.JSON.Send(ctx.Conn, map[string]interface{}{
		"msg":     "updated",
		"methods": []string{ctx.ID},
	})
}
