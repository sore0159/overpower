package main

import (
	"mule/overpower"
	"mule/overpower/db"
	"net/http"
)

func (h *Handler) SetTurnDone(w http.ResponseWriter, r *http.Request, g overpower.Game, f overpower.Faction) (isOk bool) {
	if f == nil {
		http.Error(w, "USER HAS NO FACTION IN THIS GAME", http.StatusBadRequest)
		return
	}
	if g.Turn() < 1 {
		http.Error(w, "GAME NOT IN PROGRESS", http.StatusBadRequest)
		return
	}
	facs, err := OPDB.GetFactions("gid", g.Gid())
	if my, bad := Check(err, "set turn done resource failure", "gid", g.Gid()); bad {
		return Bail(w, my)
	}
	ints, ok := GetInts(r, "turn")
	if !ok {
		http.Error(w, "MALFORMED TURN DATA", http.StatusBadRequest)
		return
	}
	turn := ints[0]
	if turn != g.Turn() {
		http.Error(w, "BAD TURN DATA", http.StatusBadRequest)
		return
	}
	done := r.FormValue("done") == "true"
	fid := f.Fid()
	var allDone bool
	if done {
		allDone = true
		for _, f := range facs {
			if f.Fid() != fid && !f.IsDone() {
				allDone = false
				break
			}
		}
	}
	if allDone {
		err = RunGameTurn(g.Gid())
		if my, bad := Check(err, "update problem in faction set done (run turn)", "faction", f); bad {
			return Bail(w, my)
		}
		return true
	} else {
		// TODO
		//f.SetDone(done)
		err = OPDB.UpdateFactions(f)
		if my, bad := Check(err, "update problem in faction set done", "faction", f); bad {
			return Bail(w, my)
		}
		return true
	}
}

func RunGameTurn(gid int) error {
	f := func(d db.DB) error {
		source := d.NewSource(gid)
		breakE, logE := overpower.RunGameTurn(source)
		if logE != nil {
			Log(logE)
		}
		if my, bad := Check(breakE, "rungameturn failed", "gid", gid); bad {
			return my
		}
		err := source.Commit()
		if my, bad := Check(err, "rungame turn commit failed", "gid", gid); bad {
			return my
		}
		return nil
	}
	return OPDB.Transact(f)
}

func StartGame(gid int, exodus bool) error {
	f := func(d db.DB) error {
		source := d.NewSource(gid)
		err := overpower.MakeGalaxy(source, exodus)
		if my, bad := Check(err, "Start game failure", "gid", gid, "exodus", exodus); bad {
			return my
		}
		err = source.Commit()
		if my, bad := Check(err, "Start game source commit failure", "gid", gid); bad {
			return my
		}
		return nil
	}
	return OPDB.Transact(f)
}
