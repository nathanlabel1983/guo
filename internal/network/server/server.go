package server

import (
	"log/slog"
	"net"
	"strconv"

	"github.com/nathanlabel1983/guo/internal/network"
)

const (
	defaultPort     = 2589
	defaulShardName = "Nathan's Shard"
)

const (
	cmdSeedMessage  = byte(0xEF) // SeedMessage is the command byte for a seed message
	cmdLoginRequest = byte(0x80)
)

// configuration holds the configuration for the server
type configuration struct {
	ShardName string // ShardName is the name of the shard
	Port      int    // Port is the port the server listens on
	PublicIP  net.IP // PublicIP is the public IP of the server
	PrivateIP net.IP // PrivateIP is the private IP of the server
}

type Server struct {
	configuration

	quitChannel chan struct{}
}

// New returns a new server with the provided options
func New(opts ...ServerOption) *Server {
	s := Default()

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Default returns a new server with default configuration
func Default() *Server {
	ip, err := network.GetPrivateIP()
	if err != nil {
		slog.Error("failed to get private IP")
		panic(err)
	}
	s := &Server{
		configuration: configuration{
			ShardName: defaulShardName,
			Port:      defaultPort,
			PrivateIP: ip,
		},
		quitChannel: make(chan struct{}),
	}
	return s
}

func (s *Server) Start() {
	slog.Info("Starting server", "IP Address", s.PrivateIP.String(), "Port", s.Port)
	s.acceptLoop()
}

func (s *Server) Stop() {
	slog.Info("Stopping server")
	close(s.quitChannel)
}

func (s *Server) acceptLoop() {
	l, err := net.Listen("tcp", s.PrivateIP.String()+":"+strconv.Itoa(s.Port))
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}
	for {
		select {
		case <-s.quitChannel:
			return
		default:
			// Accept connections
			conn, err := l.Accept()
			if err != nil {
				slog.Error("failed to accept connection", "error", err)
				continue
			}
			go s.handleConnection(conn)
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	slog.Info("Accepted Connection", "remote address", conn.RemoteAddr().String())
	cmd := make([]byte, 1)
	_, err := conn.Read(cmd)
	if err != nil {
		slog.Error("failed to read command", "error", err)
		return
	}
	switch cmd[0] {
	case cmdSeedMessage:
		data := make([]byte, 20)
		_, err := conn.Read(data)
		if err != nil {
			slog.Error("failed to read seed message", "error", err)
			return
		}
		seed := NewSeed(data)
		slog.Info("Received Seed Message", "Data", seed.String())
	}
}
