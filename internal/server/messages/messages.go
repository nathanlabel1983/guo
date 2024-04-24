package messages

import (
	"bytes"
	"io"
	"log/slog"
	"net"
)

func GetMessageFromSeedPacket(conn net.Conn) (io.Reader, error) {
	data := make([]byte, 20)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	slog.Info("received seed packet")
	return NewSeedMessage(bytes.NewReader(data)), nil
}

func GetMessageFromRequestLoginPacket(conn net.Conn) (io.Reader, error) {
	data := make([]byte, 61)
	_, err := io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}
	slog.Info("received login request packet")
	return NewRequestLoginMessage(bytes.NewReader(data)), nil
}
