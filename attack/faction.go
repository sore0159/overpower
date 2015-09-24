package attack

// Faction is the object clients should be working with
type Faction struct {
	Name        string
	OtherNames  map[int]string
	FactionID   int
	TurnDone    bool
	BuildOrders map[[4]int]Order
	View        SectorView
	TV          *TextView
}

func NewFaction() *Faction {
	return &Faction{
		OtherNames:  map[int]string{},
		BuildOrders: map[[4]int]Order{},
	}
}

type Order struct {
	Location [2]int
	Size     int
	Target   [2]int
}

func NewOrder() *Order {
	return &Order{}
}

func (f *Faction) SetOrder(amount int, source, target [2]int) {
	if amount < 1 {
		delete(f.BuildOrders, [4]int{source[0], source[1], target[0], target[1]})
		return
	}
	f.BuildOrders[[4]int{source[0], source[1], target[0], target[1]}] = Order{Location: source, Size: amount, Target: target}
}

func (f *Faction) CenterTV(center [2]int) {
	f.TV.Recenter(center)
}

func (f *Faction) NumAvail(source [2]int) (num int) {
	num = f.View.PlanetGrid[source].Launchers
	for _, o := range f.BuildOrders {
		if o.Location == source {
			num -= o.Size
		}
	}
	return
}
