package server

import (
	"testing"
)

func TestUnavailableMethod(t *testing.T) {
	s := Server{}
	h := NewMethodHandler(s)
	m := Message{Method: "test"}
	r := &TestResponse{}

	if err := h.handle(r, m); err == nil {
		t.Error("an error must be returned if method is not available")
	}
}

func TestAvailableMethod(t *testing.T) {
	s := Server{methods: make(map[string]MethodFn)}
	h := NewMethodHandler(s)
	m := Message{Method: "test"}
	r := &TestResponse{}
	b := false

	s.methods["test"] = func(ctx MethodContext) {
		b = true
	}

	if err := h.handle(r, m); err != nil {
		t.Error("an error must not be returned if method is available")
	}

	if !b {
		t.Error("method handler must be called")
	}
}
