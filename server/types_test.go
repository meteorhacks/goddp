package server

type TestResponse struct {
	_data interface{}
}

func (t *TestResponse) WriteJSON(d interface{}) error {
	t._data = d
	return nil
}

type TestRequest struct {
	Type    int
	Message []byte
	Error   error
}

func (t *TestRequest) ReadMessage() (int, []byte, error) {
	return t.Type, t.Message, t.Error
}
