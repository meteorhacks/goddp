package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/meteorhacks/goddp/utils/random"
	"github.com/meteorhacks/sockjs-go/sockjs"
	"golang.org/x/net/websocket"
)

type Server struct {
	methods map[string]MethodHandler
}

func New() Server {
	s := Server{}
	s.methods = make(map[string]MethodHandler)
	return s
}

func (s *Server) Listen(addr string) error {
	sockjs.WSHandshake = s.handshake

	sockjsOpts := sockjs.DefaultOptions
	sockjsHandler := sockjs.NewHandler("/sockjs", sockjsOpts, s.sockjsHandler)
	http.Handle("/sockjs/", sockjsHandler)

	wsServer := websocket.Server{Handler: s.wsHandler, Handshake: s.handshake}
	http.Handle("/websocket", wsServer)
	http.Handle("/sockjs/websocket", wsServer)

	return http.ListenAndServe(addr, nil)
}

func (s *Server) Method(name string, fn MethodHandler) {
	s.methods[name] = fn
}

func (s *Server) sockjsHandler(ws sockjs.Session) {
	conn := SockJSConn{ws}

	// TODO => use correct status codes
	defer ws.Close(0, "")

	for {
		msg, err := conn.ReadMessage()

		if err != nil {
			if err != io.EOF {
				fmt.Println("Error (Read Error):", err, msg)
			}

			break
		}

		s.handleMessage(&conn, &msg)
	}
}

func (s *Server) wsHandler(ws *websocket.Conn) {
	conn := WSConn{ws}
	defer ws.Close()

	for {
		msg, err := conn.ReadMessage()

		if err != nil {
			if err != io.EOF {
				fmt.Println("Error (Read Error):", err, msg)
			}

			break
		}

		s.handleMessage(&conn, &msg)
	}
}

func (s *Server) handleMessage(conn Connection, msg *Message) {
	switch msg.Msg {
	case "connect":
		s.handleConnect(conn, msg)
	case "ping":
		s.handlePing(conn, msg)
	case "method":
		s.handleMethod(conn, msg)
	default:
		fmt.Println("Error (Unknown Message Type):", msg)
		// TODO => send "error" ddp message
		break
	}
}

func (s *Server) handleConnect(conn Connection, m *Message) {
	msg := map[string]string{
		"msg":     "connected",
		"session": random.Id(17),
	}

	conn.WriteMessage(msg)
}

func (s *Server) handleMethod(conn Connection, m *Message) {
	fn, ok := s.methods[m.Method]

	if !ok {
		fmt.Println("Error: (Method '%s' Not Found)", m.Method)
		return
	}

	ctx := NewMethodContext(m, conn)
	go fn(ctx)
}

func (s *Server) handlePing(conn Connection, m *Message) {
	msg := map[string]string{
		"msg": "pong",
	}

	if m.ID != "" {
		msg["id"] = m.ID
	}

	conn.WriteMessage(msg)
}

func (s *Server) handshake(config *websocket.Config, req *http.Request) error {
	// accept all connections
	return nil
}
