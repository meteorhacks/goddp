package server

import (
	"errors"
)

type MethodContext struct {
	ID      string
	Args    []interface{}
	Res     Response
	Done    bool
	Updated bool
}

func NewMethodContext(m Message, res Response) MethodContext {
	ctx := MethodContext{}
	ctx.ID = m.ID
	ctx.Args = m.Params
	ctx.Res = res
	return ctx
}

func (ctx *MethodContext) SendResult(r interface{}) error {
	if ctx.Done {
		err := errors.New("already sent results for method")
		return err
	}

	ctx.Done = true
	return ctx.Res.WriteJSON(map[string]interface{}{
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
	return ctx.Res.WriteJSON(map[string]interface{}{
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
	return ctx.Res.WriteJSON(map[string]interface{}{
		"msg":     "updated",
		"methods": []string{ctx.ID},
	})
}
