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

type MessageReader interface {
	ToBytes() []byte
}

type Peer struct {
	connection         net.Conn
	ReceiveMessageChan chan MessageReader
	SendMessageChan    chan MessageReader

	ctx    context.Context
	cancel context.CancelFunc
}

func NewPeer(conn net.Conn) *Peer {
	ctx, cancel := context.WithCancel(context.Background())
	return &Peer{
		connection:         conn,
		ReceiveMessageChan: make(chan MessageReader, 10),
		SendMessageChan:    make(chan MessageReader, 10),

		ctx:    ctx,
		cancel: cancel,
	}
}

func (p *Peer) Close() {
	p.cancel()
	if p.connection != nil {
		slog.Info("connection closed", "remote_addr", p.connection.RemoteAddr().String())
		p.connection.Close()
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
				message, err := getMessageFromSeedPacket(p.connection)
				if err != nil {
					slog.Error("failed to read seed packet", "error", err)
					continue
				}
				p.ReceiveMessageChan <- message
			case CmdLoginRequestPacket:
				message, err := getMessageFromRequestLoginPacket(p.connection)
				if err != nil {
					slog.Error("failed to read login request packet", "error", err)
					continue
				}
				p.ReceiveMessageChan <- message
			default:
				slog.Error("unknown command", "command", cmd[0])
			}
		}
	}
}

func (p *Peer) handleSend(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			message := <-p.SendMessageChan
			_, err := p.connection.Write(message.ToBytes())
			if err != nil {
				slog.Error("failed to send message", "error", err)
				continue
			}
		}
	}
}

func getMessageFromSeedPacket(conn net.Conn) (MessageReader, error) {
	data := make([]byte, 20)
	_, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	slog.Info("received seed packet")
	return messages.NewSeedMessage(data), nil
}

func getMessageFromRequestLoginPacket(conn net.Conn) (MessageReader, error) {
	data := make([]byte, 61)
	_, err := conn.Read(data)
	if err != nil {
		return nil, err
	}
	slog.Info("received login request packet")
	return messages.NewRequestLoginMessage(data), nil
}
