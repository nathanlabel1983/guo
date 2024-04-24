package messages

import (
	"bytes"
	"io"
	"log/slog"
	"net"

	"github.com/nathanlabel1983/guo/internal/server/messages/inbound"
)

func GetMessageFromSeedPacket(conn net.Conn) (io.Reader, error) {
	data := make([]byte, 20)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	slog.Info("received seed packet")
	return inbound.NewSeedMessage(bytes.NewReader(data)), nil
}

func GetMessageFromRequestLoginPacket(conn net.Conn) (io.Reader, error) {
	data := make([]byte, 61)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	slog.Info("received login request packet")
	return inbound.NewRequestLoginMessage(bytes.NewReader(data)), nil
}
