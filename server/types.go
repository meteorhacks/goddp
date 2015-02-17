package server

type MethodHandler func([]interface{}) (interface{}, error)

// This has the all the possible fields a DDP message can have
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
