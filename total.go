package overpower

type TotallyOP struct {
	Source      Source
	Game        Game
	Factions    []Faction
	Planets     []Planet
	Orders      []Order
	Ships       []Ship
	ShipViews   []ShipView
	PlanetViews map[[2]int]PlanetView
	Reports     map[int]Report
}

func NewTotallyOP() *TotallyOP {
	return &TotallyOP{
		PlanetViews: map[[2]int]PlanetView{},
		ShipViews:   []ShipView{},
		//
	}
}

func (op *TotallyOP) SetPV(pv PlanetView) {
	op.PlanetViews[[2]int{pv.Pid(), pv.Fid()}] = pv
}

func (op *TotallyOP) AddReport(fid int, x string) {
	op.Reports[fid].AddContent(x)
}
