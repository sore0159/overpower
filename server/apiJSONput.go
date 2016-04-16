package main

import (
	"net/http"
)

// /overpower/json/RESOURCE
func (h *Handler) ApiJSONput(w http.ResponseWriter, r *http.Request) {
	if !h.LoggedIn {
		JSONUserError(w, "You must be logged in to use JSON PUT")
		return
	}
	resource, ok := h.StrAt(3)
	if !ok {
		JSONUserError(w, "NO RESOURCE SPECIFIED")
		return
	}
	if len(h.Path) > 4 {
		JSONUserError(w, "OVERSPECIFIED URL: TOO MANY ARUGMENTS")
		return
	}
	switch resource {
	case "launchorders":
		h.apiJSONputLaunchOrders(w, r)
	case "factions":
		h.apiJSONputFactions(w, r)
	case "mapviews":
		h.apiJSONputMapViews(w, r)
	default:
		JSONUserError(w, "UNKNOWN RESOURCE REQUESTED", KV{"resource", resource})
	}
}
func (h *Handler) apiJSONputFactions(w http.ResponseWriter, r *http.Request) {
	JSONUserError(w, "TODO")
}
func (h *Handler) apiJSONputLaunchOrders(w http.ResponseWriter, r *http.Request) {
	JSONUserError(w, "TODO")
}
func (h *Handler) apiJSONputMapViews(w http.ResponseWriter, r *http.Request) {
	JSONUserError(w, "TODO")
}

/*
import (
	"mule/jsend"
	"mule/overpower"
	"mule/overpower/json"
	"net/http"
)

func (h *Handler) apiJsonPUTFactions(w http.ResponseWriter, r *http.Request) {
	jF := struct {
		Password string `json:"password"`
		*json.Faction
	}{}
	err := jsend.Read(r, &jF)
	if my, bad := Check(err, "API PUT failure on data read"); bad {
		Kirk(my, w)
		return
	}
		//g, err := OPDB.GetGame("gid", jF.Gid)
		//_ = g // TODO Impliment faction creation game signup through JSON PUT
		//if err == ErrNoneFound {
			//jFail(w, 400, "params", "gid does not correspond to any existing object")
			//return
		//} else if my, bad := Check(err, "Json PUT failure on game validation check", "resource", "faction", "faction", jF); bad {
			//Kirk(my, w)
			//return
		//}
	facs, err := OPDB.GetFactions("gid", jF.Gid)
	if my, bad := Check(err, "Json PUT failure on faction validation check", "resource", "faction", "faction", jF); bad {
		Kirk(my, w)
		return
	}
	var f overpower.Faction
	allDone := true
	for _, test := range facs {
		if test.Fid() == jF.Fid {
			f = test
		} else if !test.IsDone() {
			allDone = false
		}
	}
	if f == nil {
		jFail(w, 400, "bad specification", "no faction found matching given faction data")
		return
	}
	if f.Owner() != h.User.String() {
		jFail(w, 400, "authorization", "you are not authorized for that resource")
		return
	}
	f.SetDoneBuffer(jF.DoneBuffer)
	err = OPDB.UpdateFactions(f)
	if my, bad := Check(err, "update problem in faction set done", "faction", f); bad {
		Kirk(my, w)
		return
	}
	if jF.DoneBuffer != 0 && allDone {
		err = RunGameTurn(jF.Gid)
		if my, bad := Check(err, "update problem in faction set done (run turn)", "faction", f); bad {
			Kirk(my, w)
			return
		}
	}
	jSuccess(w, nil)
}

func (h *Handler) apiJsonPUTOrders(w http.ResponseWriter, r *http.Request) {
	o := &json.Order{}
	err := jsend.Read(r, o)
	if my, bad := Check(err, "API PUT failure"); bad {
		Kirk(my, w)
		return
	}
	if o.Source == o.Target {
		jFail(w, 400, "bad specification", "source/target planets must differ")
		return
	}
	f, err := OPDB.GetFaction("gid", o.Gid, "fid", o.Fid)
	if err == ErrNoneFound {
		jFail(w, 400, "bad specification", "no faction found matching given order data")
		return
	} else if my, bad := Check(err, "Json PUT failure on faction validation check", "resource", "orders", "order", o); bad {
		Kirk(my, w)
		return
	}

	if f.Owner() != h.User.String() {
		jFail(w, 400, "authorization", "you are not authorized for that resource")
		return
	}
	planets, err := OPDB.GetPlanetsByLoc(o.Gid, o.Source, o.Target)
	if my, bad := Check(err, "Json PUT failure on planet validation check", "resource", "orders", "order", o); bad {
		Kirk(my, w)
		return
	}
	if len(planets) != 2 {
		jFail(w, 400, "bad specification", "no planets found matching given order data")
		return
	}
	var source overpower.Planet
	if planets[0].Loc() == o.Source {
		source = planets[0]
	} else if planets[1].Loc() == o.Source {
		source = planets[1]
	} else {
		jFail(w, 400, "bad specification", "bad planets found matching given order data")
		return
	}
	var powDir int
	if source.PrimaryFaction() == o.Fid {
		powDir = source.PrimaryPower()
	} else if source.SecondaryFaction() == o.Fid {
		powDir = source.SecondaryPower()
	} else {
		jFail(w, 400, "authorization", "you are not authorized for that resource")
		return
	}
	var have int
	if powDir == 0 {
		jFail(w, 400, "game mechanics", "planets must be powered to launch ships")
		return
	} else if powDir == 1 {
		have = source.Antimatter()
	} else {
		have = source.Tachyons()
	}
	using := 0
	var curOrder overpower.Order
	orders, err := OPDB.GetOrders("gid", o.Gid, "fid", o.Fid, "sourcex", o.Source[0], "sourcey", o.Source[1])
	if my, bad := Check(err, "Json PUT failure on orders validation check", "resource", "orders", "order", o); bad {
		Kirk(my, w)
		return
	}
	for _, test := range orders {
		if test.Target() != o.Target {
			using += test.Size()
		} else {
			curOrder = test
		}
	}
	if curOrder == nil && o.Size < 1 {
		jFail(w, 400, "bad specification", "size of orders must be positive on creation")
		return
	}
	if o.Size > 0 && using+o.Size > have {
		jFail(w, 400, "bad specification", "source planet has insufficient resources for order")
		return
	}
	if curOrder != nil {
		curOrder.SetSize(o.Size)
		err = OPDB.UpdateOrders(curOrder)
	} else {
		err = OPDB.MakeOrder(o.Gid, o.Fid, o.Size, o.Source, o.Target)
	}
	if my, bad := Check(err, "Json PUT failure on database entry", "resource", "orders", "item", o); bad {
		Kirk(my, w)
		return
	}
	jSuccess(w, nil)
}

func (h *Handler) apiJsonPUTMapViews(w http.ResponseWriter, r *http.Request) {
	mV := &json.MapView{}
	err := jsend.Read(r, &mV)
	if my, bad := Check(err, "API PUT failure on data read"); bad {
		Kirk(my, w)
		return
	}
	f, err := OPDB.GetFaction("gid", mV.Gid, "fid", mV.Fid)
	if err == ErrNoneFound {
		jFail(w, 400, "bad specification", "no faction found matching given mapview data")
		return
	} else if my, bad := Check(err, "Json PUT failure on mapview validation check", "resource", "mapview", "mapview", mV); bad {
		Kirk(my, w)
		return
	}
	if f.Owner() != h.User.String() {
		jFail(w, 400, "authorization", "you are not authorized for that resource")
		return
	}
	Ping("JSON PUT", mV)
	err = OPDB.UpdateMapView("center", mV.Center, "WHERE", "gid", mV.Gid, "fid", mV.Fid)
	if my, bad := Check(err, "Json PUT failure on database entry", "resource", "mapview", "item", mV); bad {
		Kirk(my, w)
		return
	}
	jSuccess(w, nil)
}
*/
