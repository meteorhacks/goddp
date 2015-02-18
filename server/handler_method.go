package server

import (
	"errors"
	"fmt"
)

type MethodHandler struct {
	server Server
}

func NewMethodHandler(s Server) Handler {
	return &MethodHandler{s}
}

func (h *MethodHandler) handle(res Response, m Message) error {
	fn, ok := h.server.methods[m.Method]

	if !ok {
		err := errors.New(fmt.Sprintf("method %s not found", m.Method))
		return err
	}

	ctx := NewMethodContext(m, res)
	fn(ctx)

	return nil
}
