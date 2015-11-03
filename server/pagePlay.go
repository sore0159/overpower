package main

import (
	"mule/hexagon"
	"mule/planetattack"
	"net/http"
)

var TPPLAY = MixTemp("frame", "titlebar", "play")

func (h *GameHandler) PlayPage() {
	if h.r.Method == "POST" {
		switch action := h.r.FormValue("action"); action {
		case "orderlaunch":
			ints, err := GetInts(h.r, "source", "target", "size")
			if err != nil {
				h.v.SetError("Invalid order form data:", err)
				goto GET
			}
			size, source, target := ints["size"], ints["source"], ints["target"]
			planets := h.g.GetPlanets(source, target)
			if planets[0] == nil {
				h.v.SetError("Can't find planet ID", source)
				goto GET
			}
			if planets[1] == nil {
				h.v.SetError("Can't find planet ID", target)
				goto GET
			}
			if planets[0].Controller != h.f.Fid {
				h.v.SetError("Faction", h.f.Name, "does not control planet", planets[0].Name, "to give orders!")
				goto GET
			}
			if planets[0].Parts < size {
				h.v.SetError("Planet", planets[0].Name, "does not have enough parts to launch a ship of size", size)
				goto GET
			}
			if source == target {
				h.v.SetError("Cannot target same planet as source!")
				goto GET
			}
			if size <= 0 {
				err := h.f.DelOrder(source, target)
				if err != nil {
					h.v.SetError("Error deleting order:", err)
					goto GET
				}
				http.Redirect(h.w, h.r, h.r.URL.Path, http.StatusFound)
				return
			}
			//orders := f.GetOrdersFrom(ints["source"])
			// TODO: check existing orders to see if we have spare parts
			orders := h.f.OrdersFor(source)
			avail := planets[0].Parts
			for _, o := range orders {
				if o.Target != target {
					avail -= o.Size
				}
			}
			if avail < size {
				h.v.SetError("Planet", planets[0].Name, "does not have enough unrequisitioned parts to launch a ship of size", size)
				goto GET
			}
			err = h.f.SetOrder(source, target, size)
			if err != nil {
				h.v.SetError("Error creating order:", err)
				goto GET
			}
		case "centerview":
			loc, err := GetInts(h.r, "x", "y")
			if err != nil {
				h.v.SetError("Invalid center coords:", err)
				goto GET
			}
			err = h.f.SetViewCenter(hexagon.Coord{loc["x"], loc["y"]})
			if err != nil {
				h.v.SetError("Recenter error:", err)
				goto GET
			}
			http.Redirect(h.w, h.r, h.r.URL.Path, http.StatusFound)
			return
		default:
			h.v.SetError("Invalid action:", action)
			goto GET
		}
	}
GET:
	m := map[string]interface{}{"userN": h.u, "game": h.g, "myfac": h.f}
	h.v.SetApp(m)
	facs, _ := h.g.Factions()
	fNames := make(map[int]string, len(facs))
	for fid, f := range facs {
		fNames[fid] = f.Name
	}
	orderL := h.f.AllOrders()
	orders := map[int][]*planetattack.Order{}
	for _, o := range orderL {
		if list, ok := orders[o.Source]; ok {
			orders[o.Source] = append(list, o)
		} else {
			orders[o.Source] = []*planetattack.Order{o}
		}
	}
	m["orders"] = orders

	center, zoom := h.f.GetView()
	m["fNames"] = fNames
	m["center"] = center
	m["zoom"] = zoom
	pvs := h.f.SortedPlanetViews(center)
	plNames := make(map[int]string, len(pvs))
	for _, pl := range pvs {
		plNames[pl.Pid] = pl.Name
		if pl.Controller == h.f.Fid {
			avail := pl.Parts
			for _, o := range orders[pl.Pid] {
				avail -= o.Size
			}
			pl.Avail = avail
		}
	}
	m["planetnames"] = plNames
	m["pvs"] = pvs
	h.v.Apply(TPPLAY, h.w)
}
