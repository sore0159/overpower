package db

import (
	"mule/hexagon"
)

type Ship struct {
	gid      int
	fid      int
	sid      int
	size     int
	launched int
	path     []hexagon.Coord
}

func NewShip() *Ship {
	return &Ship{
	//
	}
}

func (s *Ship) Gid() int {
	return s.gid
}
func (s *Ship) Fid() int {
	return s.fid
}
func (s *Ship) Sid() int {
	return s.sid
}
func (s *Ship) Size() int {
	return s.size
}
func (s *Ship) Launched() int {
	return s.launched
}
func (s *Ship) Path() []hexagon.Coord {
	return s.path
}
