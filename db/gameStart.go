package db

import (
	"mule/overpower/setup"
)

func (d DB) StartGame(gid int) (ok bool) {
	g, ok := d.GetGame(gid)
	if !ok {
		return
	}
	facs, ok := d.GetGidFactions(gid)
	if !ok {
		return
	}
	fids := make([]int, len(facs))
	for i, f := range facs {
		fids[i] = f.Fid()
	}
	planets := setup.MakeGalaxy(fids)
	return true
}
