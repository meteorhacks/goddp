package server

type SockJSSession interface {
	Recv() (string, error)
	Send(string) error
}

type Connection interface {
	ReadMessage() (Message, error)
	WriteMessage(interface{}) error
}

type MethodHandler func(MethodContext)

// This has the all the possible fields a DDP message can have
type Message struct {
	Msg     string        `json:"msg"`
	Session string        `json:"session"`
	Version string        `json:"version"`
	Support []string      `json:"support"`
	ID      string        `json:"id"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Result  string        `json:"result"`
	Methods []string      `json:"methods"`
}
