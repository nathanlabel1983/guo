package server

import (
	"context"
	"io"
	"log/slog"
	"net"

	"github.com/nathanlabel1983/guo/internal/server/messages"
	"github.com/nathanlabel1983/guo/internal/server/messages/inbound"
	"github.com/nathanlabel1983/guo/internal/server/messages/outbound"
)

const (
	CmdSeedPacket         = byte(0xEF)
	CmdLoginRequestPacket = byte(0x80)
)

// seed is the initial packet sent by the client to the server, this must be successful
// before any other packets can be sent. It will trigger the authentication process.
type seed struct {
	SeedIP string // Seed of the client, usually an IP address
	Major  uint32 // Major version of the client
	Minor  uint32 // Minor version of the client
	Rev    uint32 // Revision of the client
	Proto  uint32 // Prototype of the client
}

// Seeded returns true if the seed has been received and the client is ready to authenticate
func (s *seed) Seeded() bool {
	return s.SeedIP != ""
}

// Peer represents a connection to a client. It contains the connection, the seed information, once
// the seed packet has been received and the Peer has been authenticated, the connection is then
// handed off to the PeerHandler.
type Peer struct {
	seed

	Authenticated bool
	authenticater Authenticater

	connection         net.Conn
	ReceiveMessageChan chan io.Reader
	SendMessageChan    chan io.Reader

	ctx    context.Context
	cancel context.CancelFunc
}

type Authenticater interface {
	Authenticate(username, password string) (bool, outbound.LoginDeniedCode)
}

func NewPeer(conn net.Conn) *Peer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Peer{
		connection:         conn,
		ReceiveMessageChan: make(chan io.Reader, 10),
		SendMessageChan:    make(chan io.Reader, 10),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Peer) Close() {
	p.cancel()
	if p.connection != nil {
		slog.Info("connection closed", "remote_addr", p.connection.RemoteAddr().String())
		p.connection.Close()
		p.cancel()
	}
}

func (p *Peer) Start() {
	go p.handleReceive(p.ctx)
	go p.handleSend(p.ctx)
}

// handleReceive reads bytes from the connection and creates messages, which are then sent to the ReceiveMessageChan
func (p *Peer) handleReceive(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			cmd := make([]byte, 1)
			_, err := p.connection.Read(cmd)
			if err == io.EOF {
				p.Close()
				return
			}
			if err != nil {
				slog.Error("failed to read command", "error", err)
				p.Close()
				return
			}
			switch cmd[0] {
			case CmdSeedPacket:
				p.seedPeer()
			case CmdLoginRequestPacket:
				p.authenticatePeer()
			default:
				slog.Error("unknown command", "command", cmd[0])
			}

		}
	}
}

func (p *Peer) seedPeer() {
	seedMessage, err := messages.GetMessageFromSeedPacket(p.connection)
	if err != nil {
		panic(err)
	}
	sm, ok := seedMessage.(*inbound.SeedMessage)
	if !ok {
		panic("failed to cast seed message")
	}
	p.seed = seed{
		SeedIP: sm.IPSeed(),
		Major:  sm.Major(),
		Minor:  sm.Minor(),
		Rev:    sm.Revision(),
		Proto:  sm.Prototype(),
	}
}

func (p *Peer) authenticatePeer() {
	if !p.Seeded() {
		slog.Error("peer not seeded. Closing...")
		p.Close()
	}
	loginMessage, err := messages.GetMessageFromRequestLoginPacket(p.connection)
	if err != nil {
		slog.Error("failed to read login request packet", "error", err)
		return
	}
	lm, ok := loginMessage.(*inbound.LoginRequestMessage)
	if !ok {
		slog.Error("failed to cast login message")
		return
	}

	p.authenticater = &AuthTest{}

	authState, reason := p.authenticater.Authenticate(lm.Username(), lm.Password())
	if !authState {
		slog.Info("Authentication failed", "reason", reason)
		p.SendMessageChan <- outbound.NewLoginDeniedMessage(reason)
	}
	slog.Info("Authentication Result", "value", p.Authenticated)
}

// AuthTest is a simple Authenticater that always returns true - used for testing purposes
type AuthTest struct {
}

// Authenticate returns true
func (a *AuthTest) Authenticate(username, password string) (bool, outbound.LoginDeniedCode) {
	slog.Info("authenticating peer", "username", username, "password", password)
	return true, outbound.LoginDeniedIncorrectNamePassword
}

func (p *Peer) processPacket(pktFunc func(net.Conn) (io.Reader, error)) {
	m, err := pktFunc(p.connection)
	if err != nil {
		slog.Error("failed to read packet", "error", err)
		return
	}
	p.ReceiveMessageChan <- m
}

func (p *Peer) handleSend(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := <-p.SendMessageChan
			_, err := io.Copy(p.connection, message)
			if err != nil {
				slog.Error("failed to send message", "error", err)
				p.Close()
			}
		}
	}
}
