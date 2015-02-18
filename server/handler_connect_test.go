package server

import (
	"testing"
)

func TestHandleConnect(t *testing.T) {
	s := Server{}
	h := NewConnectHandler(s)
	m := Message{}
	r := &TestResponse{}

	h.handle(r, m)
	data := r._data.(map[string]string)

	if data["msg"] != "connected" {
		t.Error("msg field should be 'connected'")
	}

	if len(data["session"]) != 17 {
		t.Error("session field should be have 17 characters")
	}
}
