package overpower

import (
	"mule/hexagon"
)

func Travelled(sh ShipDat, turn int) (travelled []hexagon.Coord, land bool) {
	l := sh.Launched()
	if l > turn {
		return []hexagon.Coord{}, false
	}
	path := sh.Path()
	// 0 -- 10 -- 20
	start := SHIPSPEED * (turn - l)
	if start+1 > len(path) {
		return []hexagon.Coord{}, false
	}
	end := start + SHIPSPEED + 1
	if end >= len(path) {
		end = len(path)
		land = true
	}
	return path[start:end], land
}

func RadarCheck(rList, travelled []hexagon.Coord) (spotted []hexagon.Coord, spottedShip bool) {
	if len(rList) < 1 || len(travelled) < 1 {
		return nil, false
	}
	spotted = make([]hexagon.Coord, 0, len(travelled))
	for i, c := range travelled {
		for _, c2 := range rList {
			if c.StepsTo(c2) <= VISDIST {
				spotted = append(spotted, c)
				if i == len(travelled)-1 {
					spottedShip = true
				}
				break
			}
		}
	}
	return
}
