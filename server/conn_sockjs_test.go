package server

import (
	"errors"
	"testing"
)

func Test_SockJSConn_ReadMessage(t *testing.T) {
	sock := &TestSockJSSession{in: "{\"id\": \"test-id\"}"}
	conn := &SockJSConn{sock}
	msg, err := conn.ReadMessage()

	if err != nil {
		t.Error("should not return an error")
	}

	if msg.ID != "test-id" {
		t.Error("should return correct message")
	}
}

func Test_SockJSConn_ReadMessage_ReadError(t *testing.T) {
	sock := &TestSockJSSession{err: errors.New("test-err")}
	conn := &SockJSConn{sock}

	if _, err := conn.ReadMessage(); err.Error() != "test-err" {
		t.Error("should return error it read fails")
	}
}

func Test_SockJSConn_ReadMessage_UnmarshalError(t *testing.T) {
	sock := &TestSockJSSession{in: "not-json"}
	conn := &SockJSConn{sock}

	if _, err := conn.ReadMessage(); err == nil {
		t.Error("should return error if unmarshal fails")
	}
}

func Test_SockJSConn_WriteMessage(t *testing.T) {
	sock := &TestSockJSSession{}
	conn := &SockJSConn{sock}

	data := struct {
		Foo string `json:"foo"`
	}{
		Foo: "bar",
	}

	if err := conn.WriteMessage(data); err != nil {
		t.Error("should not return an error")
	}

	s, ok := conn.session.(*TestSockJSSession)
	if !ok {
		t.Error("what?")
	}

	if s.out != `{"foo":"bar"}` {
		t.Error("should send correct message")
	}
}

func Test_SockJSConn_WriteMessage_WriteError(t *testing.T) {
	sock := &TestSockJSSession{err: errors.New("test-error")}
	conn := &SockJSConn{sock}

	if err := conn.WriteMessage(""); err.Error() != "test-error" {
		t.Error("should return error if write fails")
	}
}

func Test_SockJSConn_WriteMessage_MarshalError(t *testing.T) {
	// TODO => create a Marshal error
}
