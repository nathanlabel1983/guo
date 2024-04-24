package world

import (
	"log/slog"
	"sync"
	"time"

	"github.com/nathanlabel1983/guo/internal/server"
)

type World struct {
	State
}

type State struct {
	peers     map[*server.Peer]bool
	peersMu   sync.Mutex
	timeStart time.Time
	timeNow   time.Time
}

func New() *World {
	return &World{
		State: State{
			peers:     make(map[*server.Peer]bool),
			peersMu:   sync.Mutex{},
			timeStart: time.Now(),
			timeNow:   time.Now(),
		},
	}
}

func (w *World) HandlePeer(p *server.Peer) {
	w.peersMu.Lock()
	w.peers[p] = true
	w.peersMu.Unlock()
}

func (w *World) Start() {
	go w.update()
}

// update is the main update for the world
func (w *World) update() {
	// update game ticks
	count := w.Ticks()
	for {
		w.timeNow = time.Now()
		// Print Ticks every 10 seconds
		if w.Ticks()-count >= 100 {
			slog.Info("Ticks", "ticks", w.Ticks())
			count = w.Ticks()
		}
	}
}

// Ticks returns the number of ticks that have passed since the world was created in 1/10th of a second
func (w *World) Ticks() int64 {
	return w.timeNow.Sub(w.timeStart).Nanoseconds() / 100000000
}
