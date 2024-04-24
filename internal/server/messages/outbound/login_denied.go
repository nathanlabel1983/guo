package outbound

type LoginDeniedCode byte

const (
	LoginDeniedIncorrectNamePassword LoginDeniedCode = 0x00
	LoginDeniedAlreadyInUse          LoginDeniedCode = 0x01
	LoginDeniedAccountBlocked        LoginDeniedCode = 0x02
	LoginDeniedInvalidCredentials    LoginDeniedCode = 0x03
	LoginDeniedCommunicationProblem  LoginDeniedCode = 0x04
	LoginDeniedIGRConcurrencyLimit   LoginDeniedCode = 0x05
	LoginDeniedIGRTimeLimit          LoginDeniedCode = 0x06
	LoginDeniedIGRAuthentication     LoginDeniedCode = 0x07

	CmdLoginDenied byte = 0x82
)

type LoginDeniedMessage struct {
	data []byte
}

func NewLoginDeniedMessage(reason LoginDeniedCode) *LoginDeniedMessage {
	d := make([]byte, 2)
	d[0] = CmdLoginDenied
	d[1] = byte(reason)
	return &LoginDeniedMessage{
		data: d,
	}
}

func (m *LoginDeniedMessage) Read(p []byte) (n int, err error) {
	return copy(p, m.data), nil
}

func (m *LoginDeniedMessage) Reason() LoginDeniedCode {
	return LoginDeniedCode(m.data[1])
}
