package world

import "github.com/nathanlabel1983/guo/internal/server"

type World struct {
	peers map[*server.Peer]bool
}

func New() *World {
	return &World{
		peers: make(map[*server.Peer]bool),
	}
}

func (w *World) HandlePeer(p *server.Peer) {
	w.peers[p] = true
}
