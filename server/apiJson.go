package main

import (
	"mule/jsend"
	"mule/mybad"
	"mule/overpower"
	"mule/overpower/json"
	"net/http"
)

func jFail(w http.ResponseWriter, code int, mp ...interface{}) {
	args := make(map[string]interface{}, len(mp)/2)
	for i := 0; i < len(mp)/2; i += 1 {
		args[mp[2*i].(string)] = mp[2*i+1]
	}
	err := jsend.Fail(w, code, args)
	if my, bad := Check(err, "jsend fail failure", "args", args); bad {
		Log(my)
	}
}
func jError(w http.ResponseWriter, msg string) {
	err := jsend.Error(w, msg)
	if my, bad := Check(err, "jsend error failure", "msg", msg); bad {
		Log(my)
	}
}
func jSuccess(w http.ResponseWriter, obj interface{}) {
	err := jsend.Success(w, obj)
	if my, bad := Check(err, "jsend success failure", "obj", obj); bad {
		Log(my)
	}
}

func Kirk(my *mybad.MuleError, w http.ResponseWriter) {
	Log(my)
	if my, bad := Check(jsend.Kirk(w), "api kirk failure"); bad {
		Log(my)
	}
}

// /overpower/json/
func apiJson(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		apiJsonGET(w, r)
	case "PUT":
		apiJsonPUT(w, r)
	default:
		jFail(w, 400, "method", "BAD REQUEST METHOD")
		return
	}
}

// /overpower/json/OBJ_TYPE
func apiJsonPUT(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	if !h.LoggedIn {
		jFail(w, 400, "authorization", "you are not authorized for that method")
		return
	}
	lastFull := h.LastFull()
	if lastFull < 3 {
		jFail(w, 400, "url", "no object type given")
		return
	} else if lastFull > 3 {
		jFail(w, 400, "url", "rambling url")
		return
	}
	switch h.Path[3] {
	case "orders":
		h.apiJsonPUTOrders(w, r)
	case "factions":
		h.apiJsonPUTFactions(w, r)
	default:
		jFail(w, 404, "url", "unsupported object type given")
		return
	}
}

func (h *Handler) apiJsonPUTFactions(w http.ResponseWriter, r *http.Request) {
	jF := &json.Faction{}
	err := jsend.Read(r, jF)
	if my, bad := Check(err, "API PUT failure on data read"); bad {
		Kirk(my, w)
		return
	}
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
		} else if !test.Done() {
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
	if jF.Done && allDone {
		err = RunGameTurn(jF.Gid, false)
		if my, bad := Check(err, "update problem in faction set done (run turn)", "faction", f); bad {
			Kirk(my, w)
			return
		}
	} else {
		f.SetDone(jF.Done)
		err = OPDB.UpdateFactions(f)
		if my, bad := Check(err, "update problem in faction set done", "faction", f); bad {
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
	planets, err := OPDB.GetPlanetsByPlid(o.Gid, o.Source, o.Target)
	if my, bad := Check(err, "Json PUT failure on planet validation check", "resource", "orders", "order", o); bad {
		Kirk(my, w)
		return
	}
	if len(planets) != 2 {
		jFail(w, 400, "bad specification", "no planets found matching given order data")
		return
	}
	var source overpower.Planet
	if planets[0].Pid() == o.Source {
		source = planets[0]
	} else if planets[1].Pid() == o.Source {
		source = planets[1]
	} else {
		jFail(w, 400, "bad specification", "bad planets found matching given order data")
		return
	}
	if source.Controller() != f.Fid() {
		jFail(w, 400, "authorization", "you are not authorized for that resource")
		return
	}
	have := source.Parts()
	using := 0
	var curOrder overpower.Order
	orders, err := OPDB.GetOrders("gid", o.Gid, "fid", o.Fid, "source", o.Source)
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
		jFail(w, 400, "bad specification", "source planet has insufficient parts for order")
		return
	}
	if curOrder != nil {
		curOrder.SetSize(o.Size)
		err = OPDB.UpdateOrders(curOrder)
	} else {
		err = OPDB.MakeOrder(o.Gid, o.Fid, o.Source, o.Target, o.Size)
	}
	if my, bad := Check(err, "Json PUT failure on database entry", "resource", "orders", "item", o); bad {
		Kirk(my, w)
		return
	}
	jSuccess(w, nil)
}

func apiJsonGET(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	lastFull := h.LastFull()
	if lastFull < 3 {
		jFail(w, 400, "url", "no object type given")
		return
	}
	switch h.Path[3] {
	case "games":
		h.apiJsonGETGames(w, r)
		return
	case "factions":
		h.apiJsonGETFactions(w, r)
	// overpower/json/TABLE/gid/fid/....
	case "orders", "shipviews", "reports", "planetviews":
		if !h.LoggedIn {
			jFail(w, 400, "authorization", "you are not authorized for that resource")
			return
		}
		if lastFull < 5 {
			jFail(w, 404, "url", "insuffient object specification")
			return
		}
		ints, intsOK := h.IntsAt(4, 5)
		if !intsOK {
			jFail(w, 404, "url", "badly formed object specification")
			return
		}
		searchF, err := OPDB.GetFaction("gid", ints[0], "fid", ints[1])
		if err == ErrNoneFound {
			jFail(w, 400, "authorization", "you are not authorized for that resource")
			return
		} else if my, bad := Check(err, "api get failure", "object", "faction"); bad {
			Kirk(my, w)
			return
		}
		if searchF.Owner() != h.User.String() {
			jFail(w, 400, "authorization", "you are not authorized for that resource")
			return
		}
		switch h.Path[3] {
		case "orders":
			h.apiJsonGETOrders(w, r, ints[0], ints[1])
		case "shipviews":
			h.apiJsonGETShipViews(w, r, ints[0], ints[1])
		case "planetviews":
			h.apiJsonGETPlanetViews(w, r, ints[0], ints[1])
		case "reports":
			h.apiJsonGETReports(w, r, ints[0], ints[1])
		}
		return
	default:
		jFail(w, 404, "url", "unknown object type given")
		return
	}
}

// overpower/json/orders/gid/fid/source/target
func (h *Handler) apiJsonGETOrders(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	source, ok := h.IntAt(6)
	if !ok {
		if lastFull > 5 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetOrders("gid", gid, "fid", fid)
		if my, bad := Check(err, "api json get orders failure", "gid", gid, "fid", fid); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadOrders(list)
		jSuccess(w, jsonList)
		return
	}
	target, ok := h.IntAt(7)
	if !ok {
		if lastFull > 6 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetOrders("gid", gid, "fid", fid, "source", source)
		if my, bad := Check(err, "api json get orders failure", "gid", gid, "fid", fid, "source", source); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadOrders(list)
		jSuccess(w, jsonList)
		return
	}
	if lastFull > 7 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	item, err := OPDB.GetOrder("gid", gid, "fid", fid, "source", source, "target", target)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get orders failure", "gid", gid, "fid", fid, "source", source, "target", target); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadOrder(item)
	jSuccess(w, jsonItem)
	return

}

// overpower/json/shipviews/gid/fid/turn/sid
func (h *Handler) apiJsonGETShipViews(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	turn, ok := h.IntAt(6)
	if !ok {
		if lastFull > 5 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetShipViews("gid", gid, "fid", fid)
		if my, bad := Check(err, "api json get shipviews failure", "gid", gid, "fid", fid); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadShipViews(list)
		jSuccess(w, jsonList)
		return
	}
	sid, ok := h.IntAt(7)
	if !ok {
		if lastFull > 6 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetShipViews("gid", gid, "fid", fid, "turn", turn)
		if my, bad := Check(err, "api json get shipviews failure", "gid", gid, "fid", fid, "turn", turn); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadShipViews(list)
		jSuccess(w, jsonList)
		return
	}
	if lastFull > 7 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	item, err := OPDB.GetShipView("gid", gid, "fid", fid, "turn", turn, "sid", sid)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get shipviews failure", "gid", gid, "fid", fid, "turn", turn, "sid", sid); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadShipView(item)
	jSuccess(w, jsonItem)
	return
}

// overpower/json/planetviews/gid/fid/pid
func (h *Handler) apiJsonGETPlanetViews(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	pid, ok := h.IntAt(6)
	if !ok {
		if lastFull > 5 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetPlanetViews("gid", gid, "fid", fid)
		if my, bad := Check(err, "api json get planetviews failure", "gid", gid, "fid", fid); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadPlanetViews(list)
		jSuccess(w, jsonList)
		return
	}
	if lastFull > 6 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	item, err := OPDB.GetPlanetView("gid", gid, "fid", fid, "pid", pid)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get planetviews failure", "gid", gid, "fid", fid, "pid", pid); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadPlanetView(item)
	jSuccess(w, jsonItem)
	return
}

// overpower/json/reports/gid/fid/turn
func (h *Handler) apiJsonGETReports(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	turn, ok := h.IntAt(6)
	if !ok {
		if lastFull > 5 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetReports("gid", gid, "fid", fid)
		if my, bad := Check(err, "api json get reports failure", "gid", gid, "fid", fid); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadReports(list)
		jSuccess(w, jsonList)
		return
	}
	if lastFull > 6 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	item, err := OPDB.GetReport("gid", gid, "fid", fid, "turn", turn)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get reports failure", "gid", gid, "fid", fid, "turn", turn); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadReport(item)
	jSuccess(w, jsonItem)
	return
}

// /overpower/json/games/
func (h *Handler) apiJsonGETGames(w http.ResponseWriter, r *http.Request) {
	lastFull := h.LastFull()
	if lastFull == 3 {
		games, err := OPDB.GetGames()
		if my, bad := Check(err, "api get failure", "object", "games"); bad {
			Kirk(my, w)
			return
		}
		jsonGames := json.LoadGames(games)
		jSuccess(w, jsonGames)
		return
	}
	gid, ok := h.IntAt(4)
	if !ok {
		jFail(w, 400, "params", "unreadable id for object")
		return
	}
	if lastFull == 4 {
		g, err := OPDB.GetGame("gid", gid)
		if err == ErrNoneFound {
			jFail(w, 400, "params", "id does not correspond to any existing object")
			return
		} else if my, bad := Check(err, "resource failure in json api get", "resource", "games", "gid", gid); bad {
			Kirk(my, w)
			return
		}
		jsG := json.LoadGame(g)
		jSuccess(w, jsG)
		return
	}
	jFail(w, 404, "url", "too many parameters given")
	return
}

// /overpower/json/factions/gid/fid
func (h *Handler) apiJsonGETFactions(w http.ResponseWriter, r *http.Request) {
	lastFull := h.LastFull()
	if lastFull < 5 {
		var params []interface{}
		var fid int
		if lastFull == 4 {
			gid, ok := h.IntAt(4)
			if !ok {
				jFail(w, 400, "params", "unreadable gid for object")
				return
			}
			params = []interface{}{"gid", gid}
			if h.LoggedIn {
				userF, err := OPDB.GetFaction("gid", gid, "owner", h.User.String())
				if err == ErrNoneFound {
				} else if my, bad := Check(err, "api get failure", "object", "faction"); bad {
					Log(my)
				} else {
					fid = userF.Fid()
				}
			}
		}
		list, err := OPDB.GetFactions(params...)
		if my, bad := Check(err, "api get failure", "object", "games"); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadFactions(list, fid)
		jSuccess(w, jsonList)
		return
	}
	if lastFull == 5 {
		ints, ok := h.IntsAt(4, 5)
		if !ok {
			jFail(w, 400, "params", "unreadable gid/fid for object")
			return
		}
		f, err := OPDB.GetFaction("gid", ints[0], "fid", ints[1])
		if err == ErrNoneFound {
			jFail(w, 400, "params", "id does not correspond to any existing object")
			return

		} else if my, bad := Check(err, "api get failure", "object", "faction"); bad {
			Kirk(my, w)
			return
		}
		var fid int
		if h.LoggedIn {
			userF, err := OPDB.GetFaction("gid", ints[0], "owner", h.User.String())
			if err == ErrNoneFound {
			} else if my, bad := Check(err, "api get failure", "object", "faction"); bad {
				Log(my)
			} else {
				fid = userF.Fid()
			}
		}
		jsonItem := json.LoadFaction(f, fid)
		jSuccess(w, jsonItem)
		return
	}
	jFail(w, 404, "url", "too many parameters given")
	return
}
