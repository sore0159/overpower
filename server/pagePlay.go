package main

import (
	"fmt"
	"mule/hexagon"
	"net/http"
)

var TPPLAY = MixTemp("frame", "titlebar", "play")

func (h *GameHandler) PlayPage() {
	if h.r.Method == "POST" {
		switch action := h.r.FormValue("action"); action {
		case "orderlaunch":
			ints, err := GetInts(h.r, "source", "target", "size")
			if err != nil || ints["size"] < 0 {
				h.v.SetError("Invalid order form data:", err)
				goto GET
			}
			planets := h.g.GetPlanets(ints["source"], ints["target"])
			if len(planets) < 2 {
				h.v.SetError("Can't find planets", ints["source"], ints["target"])
				goto GET
			}
			if planets[0].Controller != h.f.Fid {
				h.v.SetError("Faction", h.f.Name, "does not control planet", planets[0].Name, "to give orders!")
				goto GET
			}
			size, source, target := ints["size"], ints["source"], ints["target"]
			if planets[0].Parts < size {
				h.v.SetError("Planet", planets[0].Name, "does not have enough parts to launch a ship of size", size)
				goto GET
			}
			if size == 0 {
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
			err = h.f.SetOrder(source, target, size)
			if err != nil {
				h.v.SetError("Error creating order:", err)
				goto GET
			}
			fmt.Println(planets)
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
