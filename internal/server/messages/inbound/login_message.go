package inbound

import (
	"bytes"
	"io"
)

type LoginRequestMessage struct {
	data []byte
}

func NewRequestLoginMessage(data io.Reader) *LoginRequestMessage {
	d := make([]byte, 61)
	_, err := io.ReadFull(data, d)
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

func (s *LoginRequestMessage) Username() string {
	// remove 0x00 padding
	u := bytes.Trim(s.data[0:30], "\x00")
	return string(u)
}

func (s *LoginRequestMessage) Password() string {
	// remove 0x00 padding
	p := bytes.Trim(s.data[30:60], "\x00")
	return string(p)
}

func (s *LoginRequestMessage) NextLoginKey() byte {
	return s.data[61]
}
