package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	methods  map[string]MethodHandler
	upgrader websocket.Upgrader
}

func New() Server {
	server := Server{}
	server.methods = make(map[string]MethodHandler)
	server.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	return server
}

func (s *Server) Method(n string, h MethodHandler) {
	s.methods[n] = h
}

func (s *Server) Listen(ipPort string) {
	http.HandleFunc("/websocket", s.handler)
	http.ListenAndServe(ipPort, nil)
}

// create websocket connection from http handler and runs the websocket handler
func (s *Server) handler(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Error: could not creating websocket connection")
		return
	}

	for {
		msg, err := readMessage(ws)

		if err != nil {
			fmt.Println("Error: could not read from websocket connection")
			ws.Close()
			break
		}

		fmt.Println("Message =>", *msg)

		switch {
		case msg.Msg == "ping":
			go s.handlePing(ws, msg)
		case msg.Msg == "connect":
			go s.handleConnect(ws, msg)
		case msg.Msg == "method":
			go s.handleMethod(ws, msg)
		default:
			fmt.Println("Error: unknown ddp message", *msg)
		}
	}
}

func (s *Server) handleConnect(c *websocket.Conn, m *Message) {
	err := c.WriteJSON(map[string]string{
		"msg":     "connected",
		"session": randomId(17),
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) handlePing(c *websocket.Conn, m *Message) {
	if m.Id != "" {
		err := c.WriteJSON(map[string]string{
			"msg": "pong",
			"id":  m.Id,
		})

		if err != nil {
			fmt.Println(err)
		}
	} else {
		err := c.WriteJSON(map[string]string{
			"msg": "pong",
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s *Server) handleMethod(c *websocket.Conn, m *Message) {
	res, _ := s.methods[m.Method](m.Params)
	err := c.WriteJSON(map[string]interface{}{
		"msg":    "result",
		"id":     m.Id,
		"result": res,
	})

	if err != nil {
		fmt.Println(err)
	}

	err = c.WriteJSON(map[string]interface{}{
		"msg":     "updated",
		"methods": []string{m.Id},
	})

	if err != nil {
		fmt.Println(err)
	}
}

func readMessage(ws *websocket.Conn) (*Message, error) {
	t, str, err := ws.ReadMessage()
	msg := &Message{}

	if err != nil {
		// error reading message
		return nil, err
	}

	if t != 1 {
		// ignore binary data
		err = errors.New("Error: DDP does not supports binary streams yet.")
		return nil, err
	}

	err = json.Unmarshal(str, msg)
	return msg, nil
}
