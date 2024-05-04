package server

import "net"

type ServerOption func(*Server)

func WithShardName(name string) ServerOption {
	return func(s *Server) {
		s.ShardName = name
	}
}

func WithPort(port int) ServerOption {
	return func(s *Server) {
		s.Port = port
	}
}

func WithPrivateIP(ip string) ServerOption {
	return func(s *Server) {
		s.PrivateIP = net.ParseIP(ip)
	}
}

func WithPublicIP(ip string) ServerOption {
	return func(s *Server) {
		s.PublicIP = net.ParseIP(ip)
	}
}
