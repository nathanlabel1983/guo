package server

import (
	"context"
	"io"
	"log/slog"
	"net"

	"github.com/nathanlabel1983/guo/internal/server/messages"
)

const (
	CmdSeedPacket         = byte(0xEF)
	CmdLoginRequestPacket = byte(0x80)
)

type Seed struct {
	SeedIP string // Seed of the client, usually an IP address
	Major  uint32 // Major version of the client
	Minor  uint32 // Minor version of the client
	Rev    uint32 // Revision of the client
	Proto  uint32 // Prototype of the client
}

type Peer struct {
	Seed

	connection         net.Conn
	ReceiveMessageChan chan io.Reader
	SendMessageChan    chan io.Reader

	ctx    context.Context
	cancel context.CancelFunc
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
				continue
			}
			switch cmd[0] {
			case CmdSeedPacket:
				p.seedPeer()
				p.authenticatePeer()
			case CmdLoginRequestPacket:
				p.processPacket(messages.GetMessageFromRequestLoginPacket)
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
	sm, ok := seedMessage.(*messages.SeedMessage)
	if !ok {
		panic("failed to cast seed message")
	}
	p.Seed = Seed{
		SeedIP: sm.IPSeed(),
		Major:  sm.Major(),
		Minor:  sm.Minor(),
		Rev:    sm.Revision(),
		Proto:  sm.Prototype(),
	}
}

func (p *Peer) authenticatePeer() {

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
				continue
			}
		}
	}
}
