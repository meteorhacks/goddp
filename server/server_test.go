package server

import (
	"errors"
	"reflect"
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

func TestHandleConnect(t *testing.T) {
	s := &Server{}
	m := Message{}
	r := &TestResponse{}

	if err := handleConnect(s, r, m); err != nil {
		t.Error("connect should be handled successfully")
	}

	data := r._data.(map[string]string)
	if data["msg"] != "connected" {
		t.Error("msg field should be 'connected'")
	}

	if len(data["session"]) != 17 {
		t.Error("session field should be have 17 characters")
	}
}

func TestUnavailableMethod(t *testing.T) {
	s := &Server{}
	m := Message{Method: "test"}
	r := &TestResponse{}

	if err := handleMethod(s, r, m); err == nil {
		t.Error("an error must be returned if method is not available")
	}
}

func TestAvailableMethod(t *testing.T) {
	s := &Server{methods: make(map[string]MethodFn)}
	m := Message{Method: "test"}
	r := &TestResponse{}
	c := make(chan bool)

	s.methods["test"] = func(ctx MethodContext) {
		c <- true
	}

	if err := handleMethod(s, r, m); err != nil {
		t.Error("an error must not be returned if method is available")
	}

	// block untill method is called
	<-c
}

func TestHandlePingWithoutID(t *testing.T) {
	s := &Server{}
	m := Message{}
	r := &TestResponse{}

	if err := handlePing(s, r, m); err != nil {
		t.Error("ping should be handled successfully")
	}

	expected := map[string]string{
		"msg": "pong",
	}

	if !reflect.DeepEqual(r._data, expected) {
		t.Error("message should only have msg field")
	}
}

func TestHandlePingWithID(t *testing.T) {
	s := &Server{}
	m := Message{ID: "test-id"}
	r := &TestResponse{}

	if err := handlePing(s, r, m); err != nil {
		t.Error("ping should be handled successfully")
	}

	expected := map[string]string{
		"msg": "pong",
		"id":  "test-id",
	}

	if !reflect.DeepEqual(r._data, expected) {
		t.Error("message should have msg and ID fields")
	}
}
