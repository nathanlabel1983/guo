package messages

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
)

type SeedMessage struct {
	data []byte
}

func NewSeedMessage(data io.Reader) *SeedMessage {
	bytesRequired := 20
	d := make([]byte, bytesRequired)
	_, err := io.ReadFull(data, d)
	if err != nil {
		slog.Error("error reading seed message")
		return nil
	}
	return &SeedMessage{
		data: d,
	}
}

func (m *SeedMessage) Read(p []byte) (n int, err error) {
	return copy(p, m.data), nil
}

// IP returns the IP address of the seed message, this is the first 4 bytes as x.x.x.x
func (m *SeedMessage) IPSeed() string {
	return fmt.Sprintf("%d.%d.%d.%d", m.data[0], m.data[1], m.data[2], m.data[3])
}

func (m *SeedMessage) Major() uint32 {
	return binary.LittleEndian.Uint32(m.data[4:8])
}

func (m *SeedMessage) Minor() uint32 {
	return binary.LittleEndian.Uint32(m.data[8:12])
}

func (m *SeedMessage) Revision() uint32 {
	return binary.LittleEndian.Uint32(m.data[12:16])
}

func (m *SeedMessage) Prototype() uint32 {
	return binary.LittleEndian.Uint32(m.data[16:20])
}
