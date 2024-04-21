package messages

type LoginRequestMessage struct {
	data []byte
}

func NewRequestLoginMessage(data []byte) *LoginRequestMessage {
	return &LoginRequestMessage{
		data: data,
	}
}

func (m *LoginRequestMessage) ToBytes() []byte {
	return m.data
}
