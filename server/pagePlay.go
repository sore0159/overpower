package main

var TPPLAY = MixTemp("frame", "titlebar", "play")

func (h *GameHandler) PlayPage() {
	m := map[string]interface{}{"userN": h.u, "game": h.g, "myfac": h.f}
	h.v.SetApp(m)
	facs := h.g.Factions()
	fNames := make(map[int]string, len(facs))
	for fid, f := range facs {
		fNames[fid] = f.Name
	}
	//GET:
	center, zoom := h.f.GetView()
	m["fNames"] = fNames
	m["center"] = center
	m["zoom"] = zoom
	pvs := h.f.SortedPlanetViews(center)
	m["pvs"] = pvs
	h.v.Apply(TPPLAY, h.w)
}
