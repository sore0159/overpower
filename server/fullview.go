package main

import (
	"mule/overpower"
	"sort"
)

type FullView struct {
	Game overpower.GameDat `json:"game"`
	//	Faction       overpower.FactionDat        `json:"faction"`
	Factions      []overpower.FactionDat      `json:"factions"`
	PlanetViews   []overpower.PlanetViewDat   `json:"planetviews"`
	ShipViews     []overpower.ShipViewDat     `json:"shipviews"`
	LaunchOrders  []overpower.LaunchOrderDat  `json:"launchorders"`
	PowerOrder    overpower.PowerOrderDat     `json:"powerorder"`
	LaunchRecords []overpower.LaunchRecordDat `json:"launchrecords"`
	BattleRecords []overpower.BattleRecordDat `json:"battlerecords"`
	MapView       overpower.MapViewDat        `json:"mapview"`
	Truces        []overpower.TruceDat        `json:"truces"`
}

func (h *Handler) GetFullView(gid int) (fv *FullView, errS, errU error) {
	wGID := h.GID(gid)
	games, err := h.M.Game().SelectWhere(wGID)
	if my, bad := Check(err, "GetFullView failure on resource aquisition", "resource", "games", "gid", gid); bad {
		return nil, my, nil
	}
	if len(games) == 0 {
		return nil, nil, NewError("NO GAME FOUND")
	}
	g := games[0]
	turn := g.Turn() - 1
	if turn < 0 {
		return nil, nil, NewError("GAME HAS NOT YET BEGUN")
	}
	facs, err := h.M.Faction().SelectWhere(wGID)
	if my, bad := Check(err, "GetFullView failure on resource aquisition", "resource", "factions", "gid", gid); bad {
		return nil, my, nil
	}
	var userF overpower.FactionDat
	for _, f := range facs {
		if f.Owner() == h.User.String() {
			userF = f
			break
		}
	}
	if userF == nil {
		return nil, nil, NewError("USER HAS NO FACTION FOR GAME")
	}
	userF.SetFullJSON()

	wFID := h.FID(gid, userF.FID())
	wTURN := h.TURN(gid, userF.FID(), turn)

	plVs, err1 := h.M.PlanetView().SelectWhere(wFID)
	shVs, err2 := h.M.ShipView().SelectWhere(wTURN)
	launchOrders, err3 := h.M.LaunchOrder().SelectWhere(wFID)
	mapviews, err4 := h.M.MapView().SelectWhere(wFID)
	laRec, err5 := h.M.LaunchRecord().SelectWhere(wTURN)
	batRec, err6 := h.M.BattleRecord().SelectWhere(wTURN)
	powOrds, err7 := h.M.PowerOrder().SelectWhere(wFID)
	truces, err8 := h.M.Truce().SelectWhere(wFID)
	for i, err := range []error{err1, err2, err3, err4, err5, err6, err7, err8} {
		if my, bad := Check(err, "fill fullview failure", "index", i, "gid", gid, "fid", userF.FID(), "turn", turn); bad {
			return nil, my, nil
		}
	}
	if len(powOrds) == 0 {
		return nil, NewError("FILL FULLVIEW FAILED TO FIND POWERORDER"), nil
	}
	if len(mapviews) == 0 {
		return nil, NewError("FILL FULLVIEW FAILED TO FIND MAPVIEWS"), nil
	}

	sortLARecords(laRec)
	sortLDRecords(batRec)
	sortFactions(facs)
	fv = &FullView{
		Game: g,
		//Faction:       userF,
		Factions:      facs,
		PlanetViews:   plVs,
		ShipViews:     shVs,
		LaunchOrders:  launchOrders,
		PowerOrder:    powOrds[0],
		Truces:        truces,
		MapView:       mapviews[0],
		LaunchRecords: laRec,
		BattleRecords: batRec,
	}
	return fv, nil, nil
}

func sortLARecords(list []overpower.LaunchRecordDat) {
	sort.Sort(sortLA(list))
}
func sortLDRecords(list []overpower.BattleRecordDat) {
	sort.Sort(sortLD(list))
}
func sortFactions(list []overpower.FactionDat) {
	sort.Sort(sortFA(list))
}

type sortLA []overpower.LaunchRecordDat
type sortLD []overpower.BattleRecordDat
type sortFA []overpower.FactionDat

func (s sortLA) Len() int {
	return len(s)
}
func (s sortLA) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortLA) Less(i, j int) bool {
	sI, sJ := s[i].Source(), s[j].Source()
	if sI[0] != sJ[0] {
		return sI[0] < sJ[0]
	} else {
		return sI[1] < sJ[1]
	}
}

func (s sortLD) Len() int {
	return len(s)
}
func (s sortLD) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortLD) Less(i, j int) bool {
	sI, sJ := s[i].Loc(), s[j].Loc()
	if sI[0] != sJ[0] {
		return sI[0] < sJ[0]
	} else if sI[1] != sJ[1] {
		return sI[1] < sJ[1]
	}
	return s[i].Index() < s[j].Index()
}

func (s sortFA) Len() int {
	return len(s)
}
func (s sortFA) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s sortFA) Less(i, j int) bool {
	return s[i].FID() < s[j].FID()
}
