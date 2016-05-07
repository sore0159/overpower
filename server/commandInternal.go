package main

import (
	"mule/hexagon"
	"mule/overpower"
	"mule/overpower/models"
)

type TruceCommand struct {
	GID     int           `json:"gid"`
	FID     int           `json:"fid"`
	Loc     hexagon.Coord `json:"loc"`
	Trucees []int         `json:"trucees"`
}

func InternalSetTruce(item *TruceCommand) (errS, errU error) {
	manager := OPDB.NewManager()
	list, err := manager.Truce().Select("gid", item.GID, "fid", item.FID, "locx", item.Loc[0], "locy", item.Loc[1])
	if my, bad := Check(err, "internal set truce failure on resource aquisition", "resource", "truce", "trucecommand", item); bad {
		return my, nil
	}
	trMap := make(map[int]bool, len(item.Trucees))
	for _, fid := range item.Trucees {
		trMap[fid] = true
	}
	for _, tr := range list {
		if !trMap[tr.Trucee()] {
			tr.DELETE()
		} else {
			delete(trMap, tr.Trucee())
		}
	}
	for fid, _ := range trMap {
		newTr := &models.Truce{
			GID:    item.GID,
			FID:    item.FID,
			Loc:    item.Loc,
			Trucee: fid,
		}
		manager.CreateTruce(newTr)
	}
	err = manager.Close()
	if my, bad := Check(err, "internal set truce failure on manager close", "truce", list); bad {
		return my, nil
	}
	return nil, nil
}

func InternalSetPowerOrder(gid, fid, uppower int, loc hexagon.Coord) (errS, errU error) {
	manager := OPDB.NewManager()
	list, err := manager.PowerOrder().SelectWhere(manager.FID(gid, fid))
	if my, bad := Check(err, "internal set powerorder failure on resource aquisition", "resource", "powerorder", "gid", gid, "fid", fid); bad {
		return my, nil
	}
	if len(list) == 0 {
		return nil, NewError("no faction found for given gid/fid")
	}
	list[0].SetLoc(loc)
	list[0].SetUpPower(uppower)
	err = manager.Close()
	if my, bad := Check(err, "internal set power order failure on manager close", "power order", list[0]); bad {
		return my, nil
	}
	return nil, nil
}

func InternalSetMapCenter(gid, fid int, center hexagon.Coord) (errS, errU error) {
	manager := OPDB.NewManager()
	mvs, err := manager.MapView().SelectWhere(manager.FID(gid, fid))
	if my, bad := Check(err, "internal set mapcenter failure on resource aquisition", "resource", "mapview", "gid", gid, "fid", fid); bad {
		return my, nil
	}
	if len(mvs) == 0 {
		return nil, NewError("no faction found for given gid/fid")
	}
	mvs[0].SetCenter(center)
	err = manager.Close()
	if my, bad := Check(err, "internal set mapcenter failure on manager close", "mapview", mvs[0]); bad {
		return my, nil
	}
	return nil, nil
}

func InternalSetLaunchOrder(gid, fid, size int, source, target hexagon.Coord) (errS, errU error) {
	manager := OPDB.NewManager()
	planets, err := manager.Planet().SelectByLocs(gid, source, target)
	if my, bad := Check(err, "internal set order failure on resource aquisition", "resource", "planets", "gid", gid, "source", source, "target", target); bad {
		return my, nil
	}
	if len(planets) != 2 {
		return nil, NewError("Planets not found for given locations")
	}
	var sPl, tPl overpower.PlanetDat
	if planets[0].Loc() == source {
		sPl = planets[0]
	} else if planets[0].Loc() == target {
		tPl = planets[0]
	}
	if planets[1].Loc() == source {
		sPl = planets[1]
	} else if planets[1].Loc() == target {
		tPl = planets[1]
	}
	if tPl == nil || sPl == nil {
		return nil, NewError("Could not sort source/target planets")
	}
	var powerType, avail int
	if sPl.PrimaryFaction() == fid {
		powerType = sPl.PrimaryPower()
	} else if sPl.SecondaryFaction() == fid {
		powerType = sPl.SecondaryPower()
	} else {
		return nil, NewError("Faction not in control of source planet")
	}
	switch powerType {
	case overpower.ANTIMATTER:
		avail = sPl.Antimatter()
	case overpower.TACHYONS:
		avail = sPl.Tachyons()
	default:
		return nil, NewError("Faction is not aligned to a power type for source planet")
	}
	orders, err := manager.LaunchOrder().SelectBySource(gid, fid, source)
	if my, bad := Check(err, "internal setorder failure on resource aquisition: select order by source", "gid", gid, "fid", fid, "source", source); bad {
		return my, nil
	}
	var used int
	var o overpower.LaunchOrderDat
	for _, test := range orders {
		if test.Target() == target {
			o = test
		} else {
			used += test.Size()
		}
	}
	if size < 1 {
		if o != nil {
			o.DELETE()
			err := manager.Close()
			if my, bad := Check(err, "internal setorder failure on save order deletion", "order", o); bad {
				return my, nil
			}
		}
		return nil, nil
	}
	if size > avail-used {
		return nil, NewError("Order size is too large for available resources")
	}
	if o != nil {
		o.SetSize(size)
		err := manager.Close()
		if my, bad := Check(err, "internal setorder failure on save order update", "order", o, "size", size); bad {
			return my, nil
		}
		return nil, nil
	}
	newO := &models.LaunchOrder{
		GID:    gid,
		FID:    fid,
		Size:   size,
		Source: source,
		Target: target,
	}
	manager.CreateLaunchOrder(newO)
	err = manager.Close()
	if my, bad := Check(err, "internal setorder failure on save order create", "order", newO); bad {
		return my, nil
	}
	return nil, nil
}

func InternalSetDoneBuffer(gid, fid, buff int) (errS, errU error) {
	manager := OPDB.NewManager()
	facs, err := manager.Faction().SelectWhere(manager.GID(gid))
	if my, bad := Check(err, "command set turnbuffer failure on data retrieval", "resource", "factions", "gid", gid); bad {
		return my, nil
	}

	var f overpower.FactionDat
	allDone := buff != 0
	for _, testF := range facs {
		if testF.FID() == fid {
			f = testF
			continue
		}
		if testF.DoneBuffer() == 0 {
			allDone = false
		}
	}
	if f == nil {
		return nil, NewError("faction not found for given gid/fid")
	}

	f.SetDoneBuffer(buff)
	err = manager.Close()
	if my, bad := Check(err, "command set turnbuffer failure on updating faction", "faction", f); bad {
		return my, nil
	}
	if allDone {
		// TODO: Run multiple turns if all players have done buffers > 1
		logE, failE := OPDB.SourceTransact(f.GID(), overpower.RunGameTurn)
		if my, bad := Check(failE, "command setturn done rungame failure", "gid", f.GID()); bad {
			return my, nil
		}
		if logE != nil {
			Log(logE)
		}
	}
	return nil, nil
}
