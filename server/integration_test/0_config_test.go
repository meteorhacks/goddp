package integration

import (
	"testing"
	"time"

	"github.com/meteorhacks/goddp/server"
)

var (
	WSURL  = "http://localhost:1337/websocket"
	ORIGIN = "http://localhost:1337"
	ADDR   = "localhost:1337"
	s      server.Server
)

type MethodError struct {
	Error string `json:"error"`
}

type Message struct {
	Msg     string      `json:"msg"`
	Session string      `json:"session"`
	ID      string      `json:"id"`
	Result  float64     `json:"result"`
	Error   MethodError `json:"error"`
}

func Test_StartServer(t *testing.T) {
	s = server.New()

	s.Method("double", func(ctx server.MethodContext) {
		n, ok := ctx.Params[0].(float64)

		if !ok {
			ctx.SendError("invalid parameters")
		} else {
			ctx.SendResult(n * 2)
		}

		ctx.SendUpdated()
	})

	go s.Listen(":1337")
	time.Sleep(100 * time.Millisecond)
}
