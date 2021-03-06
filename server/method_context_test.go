package server

import (
	"reflect"
	"testing"
)

func Test_MethodContext_SendResult(t *testing.T) {
	conn := &TestConn{}
	ctx := MethodContext{ID: "test-id", Conn: conn}
	err := ctx.SendResult(100)

	expected := map[string]interface{}{
		"msg":    "result",
		"id":     "test-id",
		"result": 100,
	}

	if err != nil {
		t.Error("result should be sent successfully")
	}

	if !ctx.Done {
		t.Error("context must set that a result is sent")
	}

	if !reflect.DeepEqual(conn.out, expected) {
		t.Error("invalid response for method result")
	}
}

func Test_MethodContext_SendResult_WhenDone(t *testing.T) {
	conn := &TestConn{}
	ctx := MethodContext{ID: "test-id", Conn: conn, Done: true}
	err := ctx.SendResult(100)

	if err == nil {
		t.Error("result should be sent only once")
	}

	if conn.out != nil {
		t.Error("result should be sent only once")
	}
}

func Test_MethodContext_SendError(t *testing.T) {
	conn := &TestConn{}
	ctx := MethodContext{ID: "test-id", Conn: conn}
	err := ctx.SendError("test-error")

	expected := map[string]interface{}{
		"msg": "result",
		"id":  "test-id",
		"error": map[string]string{
			"error": "test-error",
		},
	}

	if err != nil {
		t.Error("error should be sent successfully")
	}

	if !ctx.Done {
		t.Error("context must set that a result is sent")
	}

	if !reflect.DeepEqual(conn.out, expected) {
		t.Error("invalid response for method error")
	}
}

func Test_MethodContext_SendError_WhenDone(t *testing.T) {
	conn := &TestConn{}
	ctx := MethodContext{ID: "test-id", Conn: conn, Done: true}
	err := ctx.SendError("test-error")

	if err == nil {
		t.Error("error should be sent only once")
	}

	if conn.out != nil {
		t.Error("error should be sent only once")
	}
}

func Test_MethodContext_SendUpdated(t *testing.T) {
	conn := &TestConn{}
	ctx := MethodContext{ID: "test-id", Conn: conn}
	err := ctx.SendUpdated()

	expected := map[string]interface{}{
		"msg":     "updated",
		"methods": []string{"test-id"},
	}

	if err != nil {
		t.Error("updated should be sent successfully")
	}

	if !ctx.Updated {
		t.Error("context must set that updated is sent")
	}

	if !reflect.DeepEqual(conn.out, expected) {
		t.Error("invalid response for method updated")
	}
}

func Test_MethodContext_SendUpdated_WhenDone(t *testing.T) {
	conn := &TestConn{}
	ctx := MethodContext{ID: "test-id", Conn: conn, Updated: true}
	err := ctx.SendUpdated()

	if err == nil {
		t.Error("updated message should be sent only once")
	}

	if conn.out != nil {
		t.Error("updated message should be sent only once")
	}
}
