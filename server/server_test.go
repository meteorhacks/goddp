package server

import (
	"reflect"
	"testing"
)

func Test_Server_Method(t *testing.T) {
	s := New()
	s.Method("testfn", func(MethodContext) {})
	if _, ok := s.methods["testfn"]; !ok {
		t.Error("method functionm ust be stored under methods")
	}
}

func Test_Server_HandleConnect(t *testing.T) {
	s := &Server{}
	m := &Message{}
	c := &TestConn{}

	s.handleConnect(c, m)

	data := c.out.(map[string]string)
	if data["msg"] != "connected" {
		t.Error("msg field should be 'connected'")
	}

	if len(data["session"]) != 17 {
		t.Error("session field should be have 17 characters")
	}
}

func Test_Server_HandleMethod(t *testing.T) {
	s := &Server{methods: make(map[string]MethodHandler)}
	m := &Message{Method: "test"}
	c := &TestConn{}
	ch := make(chan bool)

	s.methods["test"] = func(ctx MethodContext) {
		ch <- true
	}

	s.handleMethod(c, m)

	// block untill method is called
	<-ch
}

func Test_Server_HandlePing_WithoutID(t *testing.T) {
	s := &Server{}
	m := &Message{}
	c := &TestConn{}

	s.handlePing(c, m)

	expected := map[string]string{
		"msg": "pong",
	}

	if !reflect.DeepEqual(c.out, expected) {
		t.Error("message should only have msg field")
	}
}

func Test_Server_HandlePing_WithID(t *testing.T) {
	s := &Server{}
	m := &Message{ID: "test-id"}
	c := &TestConn{}

	s.handlePing(c, m)

	expected := map[string]string{
		"msg": "pong",
		"id":  "test-id",
	}

	if !reflect.DeepEqual(c.out, expected) {
		t.Error("message should have msg and ID fields")
	}
}
