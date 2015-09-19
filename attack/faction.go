package attack

type Faction struct {
	Name        string
	FactionID   int
	TurnDone    bool
	BuildOrders []Order
	View        SectorView
}

func NewFaction() *Faction {
	return &Faction{
		BuildOrders: []Order{},
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

func (f *Faction) AddOrder(plX, plY, size, tarX, tarY int) {
	o := Order{Location: [2]int{plX, plY}, Size: size, Target: [2]int{tarX, tarY}}
	f.BuildOrders = append(f.BuildOrders, o)
}

func (f *Faction) DropOrder(index int) {
	f.BuildOrders = append(f.BuildOrders[:index], f.BuildOrders[index+1:]...)
}
func (f *Faction) ChangeOrder(index, plX, plY, size, tarX, tarY int) {
	o := Order{Location: [2]int{plX, plY}, Size: size, Target: [2]int{tarX, tarY}}
	f.BuildOrders[index] = o
}

func (f *Faction) ToggleDone() {
	f.TurnDone = !f.TurnDone
}
