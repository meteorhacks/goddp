package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/meteorhacks/goddp/utils/random"
)

func New() Server {
	s := Server{}
	s.methods = make(map[string]MethodFn)
	s.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     checkOrigin,
	}

	return s
}

type Server struct {
	methods  map[string]MethodFn
	upgrader websocket.Upgrader
}

func (s *Server) Listen(addr string) error {
	http.HandleFunc("/websocket", s.Handler)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) Method(name string, fn MethodFn) {
	s.methods[name] = fn
}

func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// TODO => handle non-websocket requests
		return
	}

	for {
		msg, err := readMessage(ws)

		if err != nil {
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
			// TODO => send "error" ddp message
			break
		}
	}

	ws.Close()
}

func checkOrigin(r *http.Request) bool {
	return true
}

func readMessage(req Request) (Message, error) {
	t, str, err := req.ReadMessage()
	msg := Message{}

	if err != nil {
		return msg, err
	}

	if t != 1 {
		err = errors.New("DDP does not supports binary streams yet")
		return msg, err
	}

	if err := json.Unmarshal(str, &msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func handleConnect(s *Server, res Response, m Message) error {
	return res.WriteJSON(map[string]string{
		"msg":     "connected",
		"session": random.Id(17),
	})
}

func handleMethod(s *Server, res Response, m Message) error {
	fn, ok := s.methods[m.Method]

	if !ok {
		err := errors.New(fmt.Sprintf("method %s not found", m.Method))
		return err
	}

	ctx := NewMethodContext(m, res)
	go fn(ctx)

	return nil
}

func handlePing(s *Server, res Response, m Message) error {
	msg := map[string]string{
		"msg": "pong",
	}

	if m.ID != "" {
		msg["id"] = m.ID
	}

	return res.WriteJSON(msg)
}
