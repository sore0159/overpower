package overpower

func (op *TotallyOP) RunGameTurn() (ok bool) {
	names := make(map[int]string, len(op.Factions))
	for _, f := range op.Factions {
		f.SetDone(false)
		names[f.Fid()] = f.Name()
	}
	// ---- SHIPS MOVE ---- //
	//
	// ---- SHIPS LAND ---- //
	arrivals := map[int]int{}
	//
	// ---- TURN STARTS ---- //
	op.Game.IncTurn()
	turn := op.Game.Turn()
	// ---- PLANETS PRODUCE ---- //
	for _, pl := range op.Planets {
		cont := pl.Controller()
		if cont == 0 {
			continue
		}
		if inh, res := pl.Inhabitants(), pl.Resources(); inh > 0 && res > 0 {
			var prod int
			switch {
			case inh > 5 && res > 5:
				prod = 5
			case inh > res:
				prod = res
			default:
				prod = inh
			}
			pl.SetResources(res - prod)
			pl.SetParts(pl.Parts() + prod)
		}
		// ---- ARRIVALS ARRIVE ---- //
		if x := arrivals[pl.Pid()]; x > 0 {
			pl.SetInhabitants(pl.Inhabitants() + x)
		}
		// ---- PLANETS ARE SEEN ---- //
		pv, ok := op.Source.NewPlanetView(cont, turn, pl)
		if !ok {
			return false
		}
		op.SetPV(pv)
	}
	return true
}
