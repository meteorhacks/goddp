package server

import (
	"errors"
	"testing"
)

func TestAddMethod(t *testing.T) {
	s := New()
	s.Method("testfn", func(MethodContext) {})
	if _, ok := s.methods["testfn"]; !ok {
		t.Error("method functionm ust be stored under methods")
	}
}

func TestReadMessageReadError(t *testing.T) {
	req := TestRequest{Error: errors.New("test-error")}
	if _, err := readMessage(&req); err == nil {
		t.Error("an error must be returned if reading from Request fails")
	}
}

func TestReadMessageBinaryMessage(t *testing.T) {
	req := TestRequest{Type: 2}
	if _, err := readMessage(&req); err == nil {
		t.Error("an error must be returned if type is binary")
	}
}

func TestReadMessageInvalidMessage(t *testing.T) {
	str := []byte("invalid-json")
	req := TestRequest{Type: 1, Message: str}
	if _, err := readMessage(&req); err == nil {
		t.Error("an error must be returned if message is not json")
	}
}

func TestReadMessageValidMessage(t *testing.T) {
	str := []byte(`{"msg": "ping"}`)
	req := TestRequest{Type: 1, Message: str}
	msg, err := readMessage(&req)

	if err != nil {
		t.Error("message must be read successfully")
	}

	if msg.Msg != "ping" {
		t.Error("message must have correct message type")
	}
}
