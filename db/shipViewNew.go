package db

import (
	"mule/hexagon"
	"mule/overpower"
)

type ShipView struct {
	gid        int
	fid        int
	sid        int
	turn       int
	controller int
	size       int
	loc        hexagon.Coord
	locValid   bool
	trail      []hexagon.Coord
}

func NewShipView() *ShipView {
	return &ShipView{
	//
	}
}

func MakeShipView(sh overpower.Ship, fid, turn int, loc hexagon.Coord, seen []hexagon.Coord) *ShipView {
	sv := &ShipView{
		gid:        sh.Gid(),
		fid:        fid,
		sid:        sh.Sid(),
		turn:       turn,
		controller: sh.Fid(),
		size:       sh.Size(),
		//
	}
	trail := make([]hexagon.Coord, 0, len(seen))
	for _, c := range seen {
		if c == loc {
			sv.locValid = true
			sv.loc = c
		} else {
			trail = append(trail, c)
		}
	}
	sv.trail = trail
	return sv
}

func (s *ShipView) Gid() int {
	return s.gid
}
func (s *ShipView) Fid() int {
	return s.fid
}
func (s *ShipView) Sid() int {
	return s.sid
}
func (s *ShipView) Turn() int {
	return s.turn
}
func (s *ShipView) Size() int {
	return s.size
}
func (s *ShipView) Controller() int {
	return s.controller
}
func (s *ShipView) Loc() (hexagon.Coord, bool) {
	if s.locValid {
		return s.loc, true
	}
	return hexagon.Coord{0, 0}, false
}
func (s *ShipView) Trail() []hexagon.Coord {
	return s.trail
}
