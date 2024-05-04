package main

import "github.com/nathanlabel1983/guo/internal/network/server"

func main() {
	s := server.Default()
	s.Start()
}
