package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/meteorhacks/goddp/utils/random"
	"golang.org/x/net/websocket"
)

type Server struct {
	methods  map[string]MethodHandler
	wsserver websocket.Server
}

func New() Server {
	s := Server{}
	s.methods = make(map[string]MethodHandler)
	s.wsserver = websocket.Server{Handler: s.handler, Handshake: s.handshake}
	return s
}

func (s *Server) Listen(addr string) error {
	http.Handle("/websocket", s.wsserver)
	http.Handle("/sockjs/websocket", s.wsserver)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) Method(name string, fn MethodHandler) {
	s.methods[name] = fn
}

func (s *Server) handshake(config *websocket.Config, req *http.Request) error {
	// accept all connections
	return nil
}

func (s *Server) handler(ws *websocket.Conn) {
	conn := Conn{ws}
	defer ws.Close()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			if err != io.EOF {
				fmt.Println("Error (Read Error):", err, msg)
			}

			break
		}

		switch msg.Msg {
		case "connect":
			s.handleConnect(&conn, msg)
		case "ping":
			s.handlePing(&conn, msg)
		case "method":
			s.handleMethod(&conn, msg)
		default:
			fmt.Println("Error (Unknown Message Type):", msg)
			// TODO => send "error" ddp message
			break
		}
	}
}

func (s *Server) handleConnect(conn Connection, m Message) {
	msg := map[string]string{
		"msg":     "connected",
		"session": random.Id(17),
	}

	conn.WriteJSON(msg)
}

func (s *Server) handleMethod(conn Connection, m Message) {
	fn, ok := s.methods[m.Method]

	if !ok {
		fmt.Println("Error: (Method '%s' Not Found)", m.Method)
		return
	}

	ctx := NewMethodContext(m, conn)
	go fn(ctx)
}

func (s *Server) handlePing(conn Connection, m Message) {
	msg := map[string]string{
		"msg": "pong",
	}

	if m.ID != "" {
		msg["id"] = m.ID
	}

	conn.WriteJSON(msg)
}
