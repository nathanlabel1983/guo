package server

import (
	"log/slog"
	"net"
)

const (
	defaultPort      = ":2589"
	defaultShardName = "Nathan's GUO Server Emulator"
)

// PeerHandler is an interface that defines an interface for something that handles peers
// normally this is the world state.
type PeerHandler interface {
	HandlePeer(*Peer)
}

type Config struct {
	Port      string // Port is the port number that the server will listen on, for example ":2589"
	ShardName string // ShardName is the name of the shard that the server is part of
}

type Server struct {
	Config

	peers       map[*Peer]bool
	ln          net.Listener
	peerHandler PeerHandler
}

func configureDefaults(cfg *Config) {
	if len(cfg.Port) == 0 {
		cfg.Port = defaultPort
	}
	if len(cfg.ShardName) == 0 {
		cfg.ShardName = defaultShardName
	}
}

func New(cfg Config) *Server {
	configureDefaults(&cfg)
	s := &Server{
		Config: cfg,
		ln:     nil,
		peers:  make(map[*Peer]bool),
	}
	return s
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.Port)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return err
	}
	s.ln = l
	slog.Info("server started", "port", s.Port, "address", l.Addr().String())
	s.acceptLoop()
	return nil
}

func (s *Server) SetPeerHandler(ph PeerHandler) {
	s.peerHandler = ph
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err)
			continue
		}
		slog.Info("accepted connection", "remote_addr", conn.RemoteAddr().String())
		peer := NewPeer(conn)
		peer.Start()
		if s.peerHandler == nil {
			panic("peerHandler is nil")
		}
		s.peerHandler.HandlePeer(peer)
	}
}
