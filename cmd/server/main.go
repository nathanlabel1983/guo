package main

import (
	"github.com/nathanlabel1983/guo/internal/server"
	"github.com/nathanlabel1983/guo/internal/world"
)

func main() {
	s := server.New(server.Config{})
	w := world.New()
	s.SetPeerHandler(w)
	s.Start()
}
