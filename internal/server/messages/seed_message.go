package messages

import "io"

type SeedMessage struct {
	data []byte
}

func NewSeedMessage(data io.Reader) *SeedMessage {
	bytesRequired := 20
	bytesReceived := 0
	d := make([]byte, bytesRequired)
	for bytesReceived < bytesRequired {
		n, err := data.Read(d[bytesReceived:])
		if err != nil {
			return nil
		}
		bytesReceived += n
	}
	return &SeedMessage{
		data: d,
	}
}

func (m *SeedMessage) Read(p []byte) (n int, err error) {
	return copy(p, m.data), nil
}
