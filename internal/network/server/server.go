package server

import "net"

type Configuration struct {
	ShardName string
	Port      string
	PublicIP  net.IP
	PrivateIP net.IP
}

type Server struct {
	ShardName string

	Port string
}

func New(opts ...ServerOption) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}
	return s
}
