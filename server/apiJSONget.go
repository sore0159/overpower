package main

import (
	"net/http"
	"strconv"
)

// /overpower/json/...
func (h *Handler) ApiJSONget(w http.ResponseWriter, r *http.Request) {
	resource, ok := h.StrAt(3)
	if !ok {
		JSONUserError(w, "NO RESOURCE SPECIFIED")
		return
	}
	var gid int
	if resource != "games" {
		gid, ok = h.IntAt(4)
		if !ok || gid == 0 {
			JSONUserError(w, "GAME SPECIFIC RESOURCES REQUIRE GID SPECIFICATION")
			return
		}
	}
	var noAuth bool
	var getter func(...KV) (interface{}, error)
	var names []string
	switch resource {
	case "fullviews":
		fv, errS, errU := h.GetFullView(gid)
		if my, bad := Check(errS, "apiJSON failure on getfullview", "gid", gid); bad {
			JSONServerError(w, my)
			return
		}
		if errU != nil {
			JSONUserError(w, errU.Error())
			return
		}
		JSONSuccess(w, fv)
		return
	case "games":
		noAuth = true
		names = []string{"gid"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.Game().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "factions":
		noAuth = true
		names = []string{"gid", "fid"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.Faction().SelectWhere(SQLAND(args...))
			for _, fac := range obj {
				if h.LoggedIn && fac.Owner() == h.User.String() {
					fac.SetFullJSON()
				}
			}
			return obj, err
		}
	case "powerorders":
		names = []string{"gid", "fid"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.PowerOrder().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "mapviews":
		names = []string{"gid", "fid"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.MapView().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "shipviews":
		names = []string{"gid", "fid", "sid"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.ShipView().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "planetviews":
		names = []string{"gid", "fid", "locx", "locy"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.PlanetView().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "truces":
		names = []string{"gid", "fid", "locx", "locy", "trucee"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.Truce().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "battlerecords":
		names = []string{"gid", "fid", "locx", "locy", "turn", "index"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.BattleRecord().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "launchorders":
		names = []string{"gid", "fid", "sourcex", "sourcey", "targetx", "targety"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.LaunchOrder().SelectWhere(SQLAND(args...))
			return obj, err
		}
	case "launchrecords":
		names = []string{"gid", "fid", "turn", "sourcex", "sourcey", "targetx", "targety"}
		getter = func(args ...KV) (interface{}, error) {
			obj, err := h.M.LaunchRecord().SelectWhere(SQLAND(args...))
			return obj, err
		}
	default:
		JSONUserError(w, "UNKNOWN RESOURCE REQUESTED", KV{"resource", resource})
		return
	}
	if !noAuth {
		if !h.LoggedIn {
			JSONUserError(w, "NOT AUTHORIZED FOR RESOURCE")
			return
		}
		fid, ok := h.IntAt(5)
		if !ok {
			JSONUserError(w, "RESOURCE REQUESTED REQUIRES FID SPECIFICATION")
			return
		}
		facs, err := h.M.Faction().SelectWhere(h.FID(gid, fid))
		if my, bad := Check(err, "apiJSON get resource failure on faction validation check", "gid", gid, "fid", fid); bad {
			JSONServerError(w, my)
			return
		}
		if len(facs) != 1 || facs[0].Owner() != h.User.String() {
			JSONUserError(w, "NOT AUTHORIZED FOR RESOURCE")
			return
		}
	}
	h.serveJSONResource(w, getter, names)
}

// /overpower/json/RESOURCE/GID/FID/...
func (h *Handler) serveJSONResource(w http.ResponseWriter, getter func(...KV) (interface{}, error), names []string) {
	numArgs := len(h.Path) - 4
	if numArgs < 0 {
		JSONUserError(w, "NO RESOURCE SPECIFIED", KV{"arguments possible", names})
		return
	} else if numArgs > len(names) {
		JSONUserError(w, "RESOURCE OVER-SPECIFIED: TOO MANY ARGS GIVEN IN URL", KV{"arguments possible", names})
		return
	}
	args := make([]KV, 0, numArgs)
	for i, argStr := range h.Path[4:] {
		if arg, err := strconv.Atoi(argStr); err == nil {
			args = append(args, KV{names[i], arg})
		} else if i != numArgs-1 || argStr != "" {
			JSONUserError(w, "UNPARSABLE ARGUMENT IN URL", KV{"index", i + 4}, KV{"arguments possible", names})
			return
		}
	}
	obj, err := getter(args...)
	if my, bad := Check(err, "serveJSONResource getter failure", "args", args); bad {
		JSONServerError(w, my)
		return
	}
	JSONSuccess(w, obj)
}
