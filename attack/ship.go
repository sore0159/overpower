package attack

import "math/rand"

type Ship struct {
	Size      int
	FactionID int
	ILocation int
	Path      [][2]int
}

func NewShip() *Ship {
	return &Ship{Path: [][2]int{}}
}

func (c *Ship) Location() [2]int {
	return c.Path[c.ILocation]
}

func (c *Ship) Target() [2]int {
	return c.Path[len(c.Path)-1]
}

func shuffleShips(in []*Ship) (out []*Ship) {
	l := len(in)
	if l < 1 {
		return []*Ship{}
	} else if l == 1 {
		return []*Ship{in[0]}
	}
	out = make([]*Ship, l)
	for i, j := range rand.Perm(l) {
		out[i] = in[j]
	}
	return out
}

type ShipTrail struct {
	FactionID int
	Size      int
}

func NewShipTrail() *ShipTrail {
	return &ShipTrail{}
}

func (c *Ship) MakeShipTrail() *ShipTrail {
	ct := NewShipTrail()
	ct.FactionID = c.FactionID
	ct.Size = c.Size
	return ct
}
