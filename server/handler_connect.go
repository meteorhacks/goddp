package server

import (
	"github.com/meteorhacks/goddp/utils/random"
)

type ConnectHandler struct {
	server Server
}

func NewConnectHandler(s Server) Handler {
	return &ConnectHandler{s}
}

func (h *ConnectHandler) handle(res Response, m Message) error {
	return res.WriteJSON(map[string]string{
		"msg":     "connected",
		"session": random.Id(17),
	})
}
