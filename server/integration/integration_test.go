package integration

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/meteorhacks/goddp/server"
)

var (
	URL    = "http://localhost:1337/websocket"
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

func TestStartServer(t *testing.T) {
	s = server.New()

	s.Method("double", func(ctx server.MethodContext) {
		n, ok := ctx.Args[0].(float64)

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

func TestConnect(t *testing.T) {
	ws, err := newClient()
	if err != nil {
		t.Error("websocket connection failed")
	}

	defer ws.Close()

	writeMessage(ws, `{"msg": "connect", "version": "1", "support": ["1"]}`, t)
	msg := readMessage(ws, t)

	if msg.Msg != "connected" {
		t.Error("inconnect DDP message type")
	}

	if len(msg.Session) != 17 {
		t.Error("session field should be have 17 characters")
	}
}

func TestPingWithoutId(t *testing.T) {
	ws, err := newClient()
	if err != nil {
		t.Error("websocket connection failed")
	}

	defer ws.Close()

	writeMessage(ws, `{"msg": "connect", "version": "1", "support": ["1"]}`, t)
	_ = readMessage(ws, t) // ignore "connected" message

	writeMessage(ws, `{"msg": "ping"}`, t)
	msg := readMessage(ws, t)

	if msg.Msg != "pong" {
		t.Error("inconnect DDP message type")
	}
}

func TestPingWithId(t *testing.T) {
	ws, err := newClient()
	if err != nil {
		t.Error("websocket connection failed")
	}

	defer ws.Close()

	writeMessage(ws, `{"msg": "connect", "version": "1", "support": ["1"]}`, t)
	_ = readMessage(ws, t) // ignore "connected" message

	writeMessage(ws, `{"msg": "ping", "id": "test-id"}`, t)
	msg := readMessage(ws, t)

	if msg.Msg != "pong" {
		t.Error("inconnect DDP message type")
	}

	if msg.ID != "test-id" {
		t.Error("inconnect random id")
	}
}

func TestMethodResult(t *testing.T) {
	ws, err := newClient()
	if err != nil {
		t.Error("websocket connection failed")
	}

	defer ws.Close()

	writeMessage(ws, `{"msg": "connect", "version": "1", "support": ["1"]}`, t)
	_ = readMessage(ws, t) // ignore "connected" message

	writeMessage(ws, `{"msg": "method", "id": "test-id", "method": "double", "params": [2]}`, t)
	msg := readMessage(ws, t)

	if msg.Msg != "result" {
		t.Error("inconnect DDP message type")
	}

	if msg.ID != "test-id" {
		t.Error("inconnect random id")
	}

	if msg.Result != 4 {
		t.Error("inconnect method result")
	}
}

func TestMethodError(t *testing.T) {
	ws, err := newClient()
	if err != nil {
		t.Error("websocket connection failed")
	}

	defer ws.Close()

	writeMessage(ws, `{"msg": "connect", "version": "1", "support": ["1"]}`, t)
	_ = readMessage(ws, t) // ignore "connected" message

	writeMessage(ws, `{"msg": "method", "id": "test-id", "method": "double", "params": ["two"]}`, t)
	msg := readMessage(ws, t)

	if msg.Msg != "result" {
		t.Error("inconnect DDP message type")
	}

	if msg.ID != "test-id" {
		t.Error("inconnect random id")
	}

	if msg.Error.Error == "" {
		t.Error("method error should be set")
	}
}

func newClient() (*websocket.Conn, error) {
	u, _ := url.Parse(URL)
	conn, err := net.Dial("tcp", ADDR)

	if err != nil {
		return nil, err
	}

	header := http.Header{"Origin": {ORIGIN}}
	ws, _, err := websocket.NewClient(conn, u, header, 1024, 1024)
	return ws, err
}

func writeMessage(c *websocket.Conn, str string, t *testing.T) {
	w, err := c.NextWriter(websocket.TextMessage)

	if err != nil {
		t.Error("cannot create websocket write")
	}

	io.WriteString(w, str)
	w.Close()
}

func readMessage(c *websocket.Conn, t *testing.T) Message {
	op, r, err := c.NextReader()

	if op != websocket.TextMessage {
		t.Error("expecting a text message")
	}

	if err != nil {
		t.Error("cannot create reader")
	}

	str, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error("websocket read error")
	}

	msg := Message{}
	if err := json.Unmarshal(str, &msg); err != nil {
		t.Error("cannot parse websocket response")
	}

	return msg
}
