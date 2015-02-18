package server

type Request interface {
	ReadMessage() (int, []byte, error)
}

type Response interface {
	WriteJSON(interface{}) error
}

type Handler interface {
	handle(Response, Message) error
}

type MethodFn func(MethodContext)

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
