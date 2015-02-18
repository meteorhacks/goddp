package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	handlers map[string]Handler
	methods  map[string]MethodFn
	upgrader websocket.Upgrader
}

func New() Server {
	s := Server{}
	s.methods = make(map[string]MethodFn)
	s.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	s.handlers = map[string]Handler{
		"connect": NewConnectHandler(s),
		"ping":    NewPingHandler(s),
		"method":  NewMethodHandler(s),
	}

	return s
}

func (s *Server) Listen(addr string) {
	http.HandleFunc("/websocket", s.Handler)
	http.ListenAndServe(addr, nil)
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
		msg, err := readMessage(req)

		if err != nil {
			break
		}

		if h, ok := s.handlers[msg.Msg]; ok {
			go h.handle(res, msg)
		} else {
			// TODO => send "error" ddp message
			break
		}
	}

	ws.Close()
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
