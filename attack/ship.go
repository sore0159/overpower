package attack

type Ship struct {
	Size      int
	ShipID    int
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

type ShipTrail struct {
	FactionID int
	Size      int
	TrailID   int
	Count     int
	Location  [2]int
}

func NewShipTrail() *ShipTrail {
	return &ShipTrail{}
}

func (c *Ship) MakeShipTrail() *ShipTrail {
	ct := NewShipTrail()
	ct.FactionID = c.FactionID
	ct.Size = c.Size
	ct.Count = c.ILocation
	ct.Location = c.Location()
	ct.TrailID = c.ShipID
	return ct
}