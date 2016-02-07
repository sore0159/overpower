package main

import (
	enjson "encoding/json"
	"mule/overpower"
	"mule/overpower/json"
	"net/http"
)

var (
	TPCANVAS = MixTemp("frame", "titlebar", "canvas")
)

// /overpower/view/GID/play
func (h *Handler) pageCanvas(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction, facs []overpower.Faction) {
	if len(h.Path) > 5 {
		http.Redirect(w, r, h.NewPath(4), http.StatusFound)
		return
	}
	cvData, err := FillCanvasData(f)
	if my, bad := Check(err, "canvas page failure", "faction", f); bad {
		Bail(w, my)
		return
	}
	dataJson, err := enjson.Marshal(cvData)
	if my, bad := Check(err, "canvas page failure", "cvData", cvData); bad {
		Bail(w, my)
		return
	}
	m := h.DefaultApp()
	m["json"] = string(dataJson)
	m["data"] = cvData
	h.Apply(TPCANVAS, w)
}

type CanvasData struct {
	Game        *json.Game         `json:"game"`
	Faction     *json.Faction      `json:"faction"`
	Factions    []*json.Faction    `json:"factions"`
	PlanetViews []*json.PlanetView `json:"planetviews"`
	ShipViews   []*json.ShipView   `json:"shipviews"`
	Orders      []*json.Order      `json:"orders"`
	MapView     *json.MapView      `json:"mapview"`
	Reports     []*json.Report     `json:"reports"`
}

func FillCanvasData(f overpower.Faction) (*CanvasData, error) {
	gid, fid := f.Gid(), f.Fid()
	var turn int
	g, err1 := OPDB.GetGame("gid", gid)
	if g != nil {
		turn = g.Turn() - 1
	}
	facs, err2 := OPDB.GetFactions("gid", gid)
	plVs, err3 := OPDB.GetPlanetViews("gid", gid, "fid", fid)
	shVs, err4 := OPDB.GetShipViews("gid", gid, "fid", fid, "turn", turn)
	orders, err5 := OPDB.GetOrders("gid", gid, "fid", fid)
	mapview, err6 := OPDB.GetMapView("gid", gid, "fid", fid)
	reports, err7 := OPDB.GetReports("gid", gid, "fid", fid)
	for i, err := range []error{err1, err2, err3, err4, err5, err6, err7} {
		if my, bad := Check(err, "fillcanvas failure", "index", i, "gid", gid, "fid", fid); bad {
			return nil, my
		}
	}
	c := &CanvasData{
		Game:        json.LoadGame(g),
		Faction:     json.LoadFaction(f, fid),
		Factions:    json.LoadFactions(facs, fid),
		PlanetViews: json.LoadPlanetViews(plVs),
		ShipViews:   json.LoadShipViews(shVs),
		Orders:      json.LoadOrders(orders),
		MapView:     json.LoadMapView(mapview),
		Reports:     json.LoadReports(reports),
	}
	return c, nil
}
