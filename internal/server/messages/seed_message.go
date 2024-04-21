package messages

type SeedMessage struct {
	data []byte
}

func NewSeedMessage(data []byte) *SeedMessage {
	return &SeedMessage{
		data: data,
	}
}

func (m *SeedMessage) ToBytes() []byte {
	return m.data
}
