package server

import "encoding/json"

type SockJSConn struct {
	session SockJSSession
}

func (c *SockJSConn) ReadMessage() (Message, error) {
	msg := Message{}

	str, err := c.session.Recv()
	if err != nil {
		return msg, err
	}

	chars := []byte(str)
	if err = json.Unmarshal(chars, &msg); err != nil {
		return msg, err
	}

	return msg, nil
}

func (c *SockJSConn) WriteMessage(msg interface{}) error {
	chars, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	str := string(chars)
	if err = c.session.Send(str); err != nil {
		return err
	}

	return nil
}
