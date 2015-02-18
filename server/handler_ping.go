package server

type PingHandler struct {
	server Server
}

func NewPingHandler(s Server) Handler {
	return &PingHandler{s}
}

func (h *PingHandler) handle(res Response, m Message) error {
	if m.ID != "" {
		return h.withId(res, m)
	} else {
		return h.withoutId(res, m)
	}
}

func (p *PingHandler) withId(res Response, m Message) error {
	return res.WriteJSON(map[string]string{
		"msg": "pong",
		"id":  m.ID,
	})
}

func (p *PingHandler) withoutId(res Response, m Message) error {
	return res.WriteJSON(map[string]string{
		"msg": "pong",
	})
}
