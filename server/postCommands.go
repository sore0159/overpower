package main

import (
	"mule/overpower"
	"mule/overpower/models"
	"strconv"
	"strings"
)

func (h *Handler) CommandStartGame(g overpower.GameDat, facs []overpower.FactionDat, exodus bool) (errServer, errUser error) {
	if g == nil {
		return nil, NewError("USER HAS NO GAME TO START")
	}
	if g.Turn() > 0 {
		return nil, NewError("GAME ALREADY IN PROGRESS")
	}
	if len(facs) < 1 {
		return nil, NewError("GAME HAS NO PLAYERS")
	}
	f := func(source overpower.Source) (logE, failE error) {
		return nil, overpower.MakeGalaxy(source, exodus)
	}
	logE, failE := OPDB.SourceTransact(g.GID(), f)
	if my, bad := Check(failE, "command startgame failure", "gid", g.GID()); bad {
		return my, nil
	}
	if logE != nil {
		Log(logE)
	}
	return nil, nil
}

func (h *Handler) CommandSetAutos(g overpower.GameDat, dayBools [7]bool) (errServer, errUser error) {
	if g == nil {
		return nil, NewError("USER HAS NO GAME IN PROGRESS")
	}
	g.SetAutoDays(dayBools)
	err := h.M.Close()
	if my, bad := Check(err, "command setauto failure on updating game", "game", g); bad {
		return my, nil
	}
	return nil, nil
}

func (h *Handler) CommandDropGame(g overpower.GameDat) (errServer, errUser error) {
	if g == nil {
		return nil, NewError("USER HAS NO GAME IN PROGRESS")
	}
	g.DELETE()
	err := h.M.Close()
	if my, bad := Check(err, "drop game failure", "gid", g.GID()); bad {
		return my, nil
	}
	return nil, nil
}

func (h *Handler) CommandNewGame(g overpower.GameDat, password, gamename, facname, towin string) (errServer, errUser error) {
	if g != nil {
		return nil, NewError("USER ALREADY HAS GAME IN PROGRESS")
	}
	if !ValidText(gamename) {
		return nil, NewError("INVALID GAME NAME")
	}
	if password != "" && !ValidText(password) {
		return nil, NewError("INVALID GAME PASSWORD")
	}
	if facname != "" && !ValidText(facname) {
		return nil, NewError("INVALID FACTION NAME")
	}
	winI, ok := strconv.Atoi(towin)
	if ok != nil || winI < 2 {
		return nil, NewError("INVALID GAME WIN THRESHOLD")
	}
	newG := &models.Game{
		Owner: h.User.String(),
		Name:  gamename,
		ToWin: winI,
	}
	if password != "" {
		newG.Password.Valid = true
		newG.Password.String = password
	}
	h.M.CreateGame(newG)
	err := h.M.Close()
	if my, bad := Check(err, "make game failure", "user", h.User, "gamename", gamename, "password", password, "towin", winI); bad {
		return my, nil
	}
	if facname == "" {
		return nil, nil
	}
	newF := &models.Faction{
		GID:   newG.GID,
		Owner: h.User.String(),
		Name:  facname,
	}
	h.M.CreateFaction(newF)
	err = h.M.Close()
	if my, bad := Check(err, "make faction failure", "user", h.User, "facname", facname, "gid", newG.GID); bad {
		return my, nil
	}
	return nil, nil
}

func (h *Handler) CommandDropFaction(g overpower.GameDat, f overpower.FactionDat) (errServer, errUser error) {
	if g.Turn() > 0 {
		return nil, NewError("GAME IN PROGRESS")
	}
	if f == nil {
		return nil, NewError("USER HAS NO FACTION FOR THIS GAME")
	}
	f.DELETE()
	err := h.M.Close()
	if my, bad := Check(err, "command drop faction failure", "game", g, "faction", f); bad {
		return my, nil
	}
	return nil, nil
}

func (h *Handler) CommandNewFaction(g overpower.GameDat, facs []overpower.FactionDat, f overpower.FactionDat, password, facname string) (errServer, errUser error) {
	if g.Turn() > 0 {
		return nil, NewError("GAME IN PROGRESS")
	}
	if f != nil {
		return nil, NewError("USER ALREADY HAS FACTION FOR THIS GAME")
	}
	if g.HasPassword() {
		if !ValidText(password) || !g.IsPassword(password) {
			return nil, NewError("BAD PASSWORD")
		}
	}
	if !ValidText(facname) {
		return nil, NewError("BAD FACTION NAME")
	}
	lwFName := strings.ToLower(facname)
	for _, f := range facs {
		if strings.ToLower(f.Name()) == lwFName {
			return nil, NewError("FACTION NAME ALREADY IN USE FOR THIS GAME")
		}
	}
	newF := &models.Faction{
		GID:   g.GID(),
		Owner: h.User.String(),
		Name:  facname,
	}
	h.M.CreateFaction(newF)
	err := h.M.Close()
	if my, bad := Check(err, "data creation error", "type", "faction", "gid", g.GID(), "user", h.User, "facname", facname); bad {
		return my, nil
	}
	return nil, nil
}

func (h *Handler) CommandQuitGame(g overpower.GameDat, f overpower.FactionDat, turnStr string) (errServer, errUser error) {
	turnI, err := strconv.Atoi(turnStr)
	if err != nil || turnI != g.Turn() {
		return nil, NewError("FORM SUBMISSION TURN DOES NOT MATCH GAME TURN")
	}
	f.DELETE()
	return h.M.Close(), nil
}

func (h *Handler) CommandSetDoneBuffer(g overpower.GameDat, facs []overpower.FactionDat, f overpower.FactionDat, turnStr, buffStr string) (errServer, errUser error) {
	if g.Turn() < 1 {
		return nil, NewError("GAME HAS NOT YET BEGUN")
	}
	turnI, err := strconv.Atoi(turnStr)
	if err != nil || turnI != g.Turn() {
		return nil, NewError("FORM SUBMISSION TURN DOES NOT MATCH GAME TURN")
	}
	buffI, err := strconv.Atoi(buffStr)
	if err != nil {
		return nil, NewError("UNPARSABLE TURN BUFFER VALUE")
	}
	if buffI < 0 {
		buffI = -1
	}
	curBuff := f.DoneBuffer()
	if buffI == curBuff {
		return nil, nil
	}
	f.SetDoneBuffer(buffI)
	err = h.M.Close()
	if my, bad := Check(err, "command set turnbuffer failure on updating faction", "faction", f); bad {
		return my, nil
	}
	var allDone bool
	if curBuff == 0 {
		allDone = true
		for _, testF := range facs {
			if testF == nil {
				continue
			}
			if testF.FID() != f.FID() && testF.DoneBuffer() == 0 {
				allDone = false
				break
			}
		}
	}

	if allDone {
		// TODO: Run multiple turns if all players have done buffers > 1
		fRun := func(source overpower.Source) (logE, failE error) {
			return overpower.RunGameTurn(source)
		}
		logE, failE := OPDB.SourceTransact(g.GID(), fRun)
		if my, bad := Check(failE, "command setturn done rungame failure", "gid", g.GID()); bad {
			return my, nil
		}
		if logE != nil {
			Log(logE)
		}
	}
	return nil, nil
}

func (h *Handler) CommandForceTurn(g overpower.GameDat, turnStr string) (errServer, errUser error) {
	if g == nil {
		return nil, NewError("USER HAS NO GAME TO PROGRESS")
	}
	if g.Turn() < 1 {
		return nil, NewError("GAME HAS NOT YET BEGUN")
	}
	turnI, err := strconv.Atoi(turnStr)
	if err != nil || turnI != g.Turn() {
		return nil, NewError("FORM SUBMISSION TURN DOES NOT MATCH GAME TURN")
	}
	logE, failE := OPDB.SourceTransact(g.GID(), overpower.RunGameTurn)
	if my, bad := Check(failE, "failure on running turn", "gid", g.GID()); bad {
		return my, nil
	}
	if logE != nil {
		Log(logE)
	}
	return nil, nil
}
