package server

type TestConn struct {
	in  Message
	out interface{}
	err error
}

func (c *TestConn) ReadMessage() (Message, error) {
	return c.in, c.err
}

func (c *TestConn) WriteMessage(msg interface{}) error {
	c.out = msg
	return c.err
}

type TestSockJSSession struct {
	in  string
	out string
	err error
}

func (s *TestSockJSSession) Recv() (string, error) {
	return s.in, s.err
}

func (s *TestSockJSSession) Send(msg string) error {
	s.out = msg
	return s.err
}
