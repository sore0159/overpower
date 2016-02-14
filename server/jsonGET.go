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
	case "orders", "shipviews", "launchrecords", "landingrecords", "planetviews", "fullviews":
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
		case "launchrecords":
			h.apiJsonGETLaunchRecords(w, r, ints[0], ints[1])
		case "landingrecords":
			h.apiJsonGETLandingRecords(w, r, ints[0], ints[1])
		case "fullviews":
			h.apiJsonGETFullViews(w, r, searchF)
		}
		return
	default:
		jFail(w, 404, "url", "unknown object type given")
		return
	}
}

// overpower/json/launchrecords/gid/fid/turn/sourcex/sourcey/targetx/targety
func (h *Handler) apiJsonGETLaunchRecords(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	if len(h.Path) > 11 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	args := []interface{}{"gid", gid, "fid", fid}
	argCheck := []string{"turn", "sourcex", "sourcey", "targetx", "targety"}
	for i, key := range argCheck {
		val, ok := h.IntAt(6 + i)
		if !ok && lastFull > 5+i {
			jFail(w, 400, "url", "bad args for object specification")
			return
		} else if !ok {
			list, err := OPDB.GetLaunchRecords(args...)
			if my, bad := Check(err, "api json get launch records failure", "args", args); bad {
				Kirk(my, w)
				return
			}
			jsonList := json.LoadLaunchRecords(list)
			jSuccess(w, jsonList)
			return
		} else {
			args = append(args, key, val)
		}
	}
	item, err := OPDB.GetLaunchRecord(args...)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get launchrecords failure", "args", args); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadLaunchRecord(item)
	jSuccess(w, jsonItem)
}

// overpower/json/landingrecords/gid/fid/turn/index
func (h *Handler) apiJsonGETLandingRecords(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	if len(h.Path) > 8 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	args := []interface{}{"gid", gid, "fid", fid}
	argCheck := []string{"turn", "index"}
	for i, key := range argCheck {
		val, ok := h.IntAt(6 + i)
		if !ok && lastFull > 5+i {
			jFail(w, 400, "url", "bad args for object specification")
			return
		} else if !ok {
			list, err := OPDB.GetLandingRecords(args...)
			if my, bad := Check(err, "api json get landing records failure", "args", args); bad {
				Kirk(my, w)
				return
			}
			jsonList := json.LoadLandingRecords(list)
			jSuccess(w, jsonList)
			return
		} else {
			args = append(args, key, val)
		}
	}
	item, err := OPDB.GetLandingRecord(args...)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get landingrecords failure", "args", args); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadLandingRecord(item)
	jSuccess(w, jsonItem)
}

// overpower/json/fullviews/gid/fid
func (h *Handler) apiJsonGETFullViews(w http.ResponseWriter, r *http.Request, f overpower.Faction) {
	if len(h.Path) > 6 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	fv, err := FillFullView(f)
	if my, bad := Check(err, "json API get failure", "resource", "fullview", "faction", f); bad {
		Kirk(my, w)
		return
	}
	jSuccess(w, fv)
}

// overpower/json/orders/gid/fid/sourcex/sourcey/targetx/targety
func (h *Handler) apiJsonGETOrders(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	sourcex, ok := h.IntAt(6)
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
	sourcey, ok := h.IntAt(7)
	if !ok {
		jFail(w, 400, "url", "bad source-y arg for object specification")
		return
	}

	targetx, ok := h.IntAt(8)
	if !ok {
		if lastFull > 7 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetOrders("gid", gid, "fid", fid, "sourcex", sourcex, "sourcey", sourcey)
		if my, bad := Check(err, "api json get orders failure", "gid", gid, "fid", fid, "sourcex", sourcex, "sourcey", sourcey); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadOrders(list)
		jSuccess(w, jsonList)
		return
	}
	targety, ok := h.IntAt(9)
	if !ok {
		jFail(w, 400, "url", "bad target-y arg for object specification")
		return
	}
	if lastFull > 9 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	item, err := OPDB.GetOrder("gid", gid, "fid", fid, "sourcex", sourcex, "sourcey", sourcey, "targetx", targetx, "targety", targety)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get orders failure", "gid", gid, "fid", fid, "sourcex", sourcex, "sourcey", sourcey, "targetx", targetx, "targety", targety); bad {
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

// overpower/json/planetviews/gid/fid/x/y
func (h *Handler) apiJsonGETPlanetViews(w http.ResponseWriter, r *http.Request, gid, fid int) {
	lastFull := h.LastFull()
	plX, ok := h.IntAt(6)
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
	plY, ok := h.IntAt(7)
	if !ok {
		if lastFull > 6 {
			jFail(w, 400, "url", "bad args for object specification")
			return
		}
		list, err := OPDB.GetPlanetViews("gid", gid, "fid", fid, "locx", plX)
		if my, bad := Check(err, "api json get planetviews failure", "gid", gid, "fid", fid, "locx", plX); bad {
			Kirk(my, w)
			return
		}
		jsonList := json.LoadPlanetViews(list)
		jSuccess(w, jsonList)
		return
	}
	if lastFull > 7 {
		jFail(w, 400, "url", "too many args for object specification")
		return
	}
	item, err := OPDB.GetPlanetView("gid", gid, "fid", fid, "locx", plX, "locy", plY)
	if err == ErrNoneFound {
		jFail(w, 400, "params", "id does not correspond to any existing object")
		return
	} else if my, bad := Check(err, "api json get planetviews failure", "gid", gid, "fid", fid, "locx", plX, "locy", plY); bad {
		Kirk(my, w)
		return
	}
	jsonItem := json.LoadPlanetView(item)
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

type FullView struct {
	Game           *json.Game            `json:"game"`
	Faction        *json.Faction         `json:"faction"`
	Factions       []*json.Faction       `json:"factions"`
	PlanetViews    []*json.PlanetView    `json:"planetviews"`
	ShipViews      []*json.ShipView      `json:"shipviews"`
	Orders         []*json.Order         `json:"orders"`
	LaunchRecords  []*json.LaunchRecord  `json:"launchrecords"`
	LandingRecords []*json.LandingRecord `json:"landingrecords"`
	MapView        *json.MapView         `json:"mapview"`
}

func FillFullView(f overpower.Faction) (*FullView, error) {
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
	laRep, err7 := OPDB.GetLaunchRecords("gid", gid, "fid", fid, "turn", turn)
	ldRep, err8 := OPDB.GetLandingRecords("gid", gid, "fid", fid, "turn", turn)
	for i, err := range []error{err1, err2, err3, err4, err5, err6, err7, err8} {
		if my, bad := Check(err, "fill fullview failure", "index", i, "gid", gid, "fid", fid); bad {
			return nil, my
		}
	}
	fv := &FullView{
		Game:           json.LoadGame(g),
		Faction:        json.LoadFaction(f, fid),
		Factions:       json.LoadFactions(facs, fid),
		PlanetViews:    json.LoadPlanetViews(plVs),
		ShipViews:      json.LoadShipViews(shVs),
		Orders:         json.LoadOrders(orders),
		MapView:        json.LoadMapView(mapview),
		LaunchRecords:  json.LoadLaunchRecords(laRep),
		LandingRecords: json.LoadLandingRecords(ldRep),
	}
	return fv, nil
}
