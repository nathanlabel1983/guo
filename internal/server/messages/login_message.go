package messages

import "io"

type LoginRequestMessage struct {
	data []byte
}

func NewRequestLoginMessage(data io.Reader) *LoginRequestMessage {
	d := make([]byte, 61)
	_, err := data.Read(d)
	if err != nil {
		return nil
	}
	return &LoginRequestMessage{
		data: d,
	}
}

func (s *LoginRequestMessage) Read(p []byte) (n int, err error) {
	return copy(p, s.data), nil
}
