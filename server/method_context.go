package server

import "errors"

type MethodContext struct {
	ID      string
	Params  []interface{}
	Conn    Connection
	Done    bool
	Updated bool
}

func NewMethodContext(m Message, conn Connection) MethodContext {
	ctx := MethodContext{}
	ctx.ID = m.ID
	ctx.Params = m.Params
	ctx.Conn = conn
	return ctx
}

func (ctx *MethodContext) SendResult(result interface{}) error {
	if ctx.Done {
		err := errors.New("results already sent")
		return err
	}

	ctx.Done = true
	msg := map[string]interface{}{
		"msg":    "result",
		"id":     ctx.ID,
		"result": result,
	}

	return ctx.Conn.WriteJSON(msg)
}

func (ctx *MethodContext) SendError(e string) error {
	if ctx.Done {
		err := errors.New("already sent results for method")
		return err
	}

	ctx.Done = true
	msg := map[string]interface{}{
		"msg": "result",
		"id":  ctx.ID,
		"error": map[string]string{
			"error": e,
		},
	}

	return ctx.Conn.WriteJSON(msg)
}

func (ctx *MethodContext) SendUpdated() error {
	if ctx.Updated {
		err := errors.New("already sent updated for method")
		return err
	}

	ctx.Updated = true
	msg := map[string]interface{}{
		"msg":     "updated",
		"methods": []string{ctx.ID},
	}

	return ctx.Conn.WriteJSON(msg)
}
