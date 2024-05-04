package server

import (
	"encoding/binary"
	"fmt"
)

type Seed struct {
	data []byte
}

func NewSeed(data []byte) *Seed {
	s := &Seed{
		data: make([]byte, 20),
	}
	copy(s.data, data)
	return s
}

func (s *Seed) String() string {
	return fmt.Sprintf("Seed: %v, Major: %v, Minor: %v, Revision: %v, Prototype: %v", s.Seed(), s.Major(), s.Minor(), s.Revision(), s.Prototype())
}

func (s *Seed) Seed() string {
	return fmt.Sprintf("%v.%v.%v.%v", s.data[0], s.data[1], s.data[2], s.data[3])
}

func (s *Seed) Major() uint32 {
	fmt.Printf("s.data[4:8] = %v\n", s.data[4:8])
	return binary.BigEndian.Uint32(s.data[4:8])
}

func (s *Seed) Minor() uint32 {
	return binary.BigEndian.Uint32(s.data[8:12])
}

func (s *Seed) Revision() uint32 {
	return binary.BigEndian.Uint32(s.data[12:16])
}

func (s *Seed) Prototype() uint32 {
	return binary.BigEndian.Uint32(s.data[16:20])
}

func (s *Seed) Write(p []byte) (n int, err error) {
	if len(p) != 20 {
		return 0, fmt.Errorf("buffer not 20")
	}
	copy(s.data, p)
	return 20, nil
}
