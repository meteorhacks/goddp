package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/meteorhacks/goddp/utils/random"
	"golang.org/x/net/websocket"
)

func New() Server {
	s := Server{}
	s.wsserver = websocket.Server{Handler: s.wsHandler, Handshake: handshake}
	s.methods = make(map[string]MethodFn)
	return s
}

type Server struct {
	methods  map[string]MethodFn
	wsserver websocket.Server
}

func (s *Server) Listen(addr string) error {
	http.Handle("/websocket", s.wsserver)
	http.Handle("/sockjs/websocket", s.wsserver)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) Method(name string, fn MethodFn) {
	s.methods[name] = fn
}

func (s *Server) wsHandler(ws *websocket.Conn) {
	for {
		var msg Message
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			if err != io.EOF {
				fmt.Println("Read Error: ", err)
			}

			break
		}

		switch msg.Msg {
		case "connect":
			handleConnect(s, ws, msg)
		case "ping":
			handlePing(s, ws, msg)
		case "method":
			handleMethod(s, ws, msg)
		default:
			fmt.Println("Error: unknown message type", msg)
			// TODO => send "error" ddp message
			break
		}
	}

	ws.Close()
}

func handshake(config *websocket.Config, req *http.Request) error {
	// accept all connections
	return nil
}

func handleConnect(s *Server, ws *websocket.Conn, m Message) error {
	return websocket.JSON.Send(ws, map[string]string{
		"msg":     "connected",
		"session": random.Id(17),
	})
}

func handleMethod(s *Server, ws *websocket.Conn, m Message) error {
	fn, ok := s.methods[m.Method]

	if !ok {
		err := errors.New(fmt.Sprintf("method %s not found", m.Method))
		return err
	}

	ctx := NewMethodContext(m, ws)
	go fn(ctx)

	return nil
}

func handlePing(s *Server, ws *websocket.Conn, m Message) error {
	msg := map[string]string{
		"msg": "pong",
	}

	if m.ID != "" {
		msg["id"] = m.ID
	}

	return websocket.JSON.Send(ws, msg)
}
