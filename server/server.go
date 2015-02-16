package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/gorilla/websocket"
)

type MethodHandler func([]interface{}) (interface{}, error)

type Server struct {
	ws      *websocket.Conn
	methods map[string]MethodHandler
}

type Message struct {
	Msg     string        `json:"msg"`
	Session string        `json:"session"`
	Version string        `json:"version"`
	Support []string      `json:"support"`
	Id      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Result  string        `json:"result"`
	Methods []string      `json:"methods"`
}

// type ConnectMessage struct {
// 	Msg     string   `json:"msg"`
// 	Session string   `json:"session"`
// 	Version string   `json:"version"`
// 	Support []string `json:"support"`
// }

// type ConnectedMessage struct {
// 	Msg     string `json:"msg"`
// 	Session string `json:"session"`
// }

// type PingWithIdMessage struct {
// 	Msg string `json:"msg"`
// 	Id  string `json:"id"`
// }

// type PingWithoutIdMessage struct {
// 	Msg string `json:"msg"`
// }

// type PongWithIdMessage struct {
// 	Msg string `json:"msg"`
// 	Id  string `json:"id"`
// }

// type PongWithoutIdMessage struct {
// 	Msg string `json:"msg"`
// }

// type MethodMessage struct {
// 	Msg    string        `json:"msg"`
// 	Id     string        `json:"id"`
// 	Method string        `json:"method"`
// 	Params []interface{} `json:"params"`
// }

// type ResultWithErrorMessage struct {
// 	Msg   string `json:"msg"`
// 	Id    string `json:"id"`
// 	Error Error  `json:"error"`
// }

// type ResultWithoutErrorMessage struct {
// 	Msg    string `json:"msg"`
// 	Id     string `json:"id"`
// 	Result string `json:"result"`
// }

// type UpdatedMessage struct {
// 	Msg     string   `json:"msg"`
// 	Methods []string `json:"methods"`
// }

// type Error struct {
// 	error   string
// 	reason  string
// 	details string
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() Server {
	server := Server{}
	server.methods = make(map[string]MethodHandler)
	return server
}

func (s *Server) Method(n string, h MethodHandler) {
	s.methods[n] = h
}

// create websocket connection from http handler and runs the websocket handler
func (s *Server) Handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	s.ws = ws

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

		fmt.Println("Message =>", msg)

		switch {
		case msg.Msg == "ping":
			go s.HandlePing(ws, &msg)
		case msg.Msg == "connect":
			go s.HandleConnect(ws, &msg)
		case msg.Msg == "method":
			go s.HandleMethod(ws, &msg)
		default:
			fmt.Println("Error: unknown ddp message", msg)
		}
	}
}

func (s *Server) HandleConnect(c *websocket.Conn, m *Message) {
	err := c.WriteJSON(map[string]string{
		"msg":     "connected",
		"session": randId(17),
	})

	if err != nil {
		fmt.Println(err)
	}
}

func (s *Server) HandlePing(c *websocket.Conn, m *Message) {
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

func (s *Server) HandleMethod(c *websocket.Conn, m *Message) {
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

func readMessage(ws *websocket.Conn) (Message, error) {
	t, str, err := ws.ReadMessage()
	msg := Message{}

	if err != nil {
		// error reading message
		return msg, err
	}

	if t != 1 {
		// ignore binary data
		err = errors.New("Error: websocket data is binary")
		return msg, err
	}

	err = json.Unmarshal(str, &msg)
	return msg, err
}

// TODO: move this to another package
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randId(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
