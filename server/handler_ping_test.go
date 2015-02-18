package server

import (
	"reflect"
	"testing"
)

func TestHandlePingWithoutID(t *testing.T) {
	s := Server{}
	h := NewPingHandler(s)
	m := Message{}
	r := &TestResponse{}
	h.handle(r, m)

	expected := map[string]string{
		"msg": "pong",
	}

	if !reflect.DeepEqual(r._data, expected) {
		t.Error("message should only have msg field")
	}
}

func TestHandlePingWithID(t *testing.T) {
	s := Server{}
	h := NewPingHandler(s)
	m := Message{ID: "test-id"}
	r := &TestResponse{}
	h.handle(r, m)

	expected := map[string]string{
		"msg": "pong",
		"id":  "test-id",
	}

	if !reflect.DeepEqual(r._data, expected) {
		t.Error("message should have msg and ID fields")
	}
}
