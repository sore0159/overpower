package overpower

import (
//	"mule/hexagon"
)

func Battle(source Source, pl PlanetDat, sh ShipDat, turn int, truces map[[2]int]TruceDat) {
	// if truce is broken (by trucee):
	// delete(truces, [2]int{trucer, trucee}
	// source.DropTrouce(pl.Loc, trucer, trucee)
	prFid := pl.PrimaryFaction()
	prPr := pl.PrimaryPresence()
	prPW := pl.PrimaryPower()
	seFid := pl.SecondaryFaction()
	sePr := pl.SecondaryPresence()
	sePW := pl.SecondaryPower()
	var shFid, shSize int
	if sh != nil {
		shFid = sh.FID()
		shSize = sh.Size()
		if shSize < 1 {
			return
		}
	}
	defer AllSee(source, pl, prFid, seFid, shFid, turn)
	var betrayals [][2]int
	defer AllBattleRecords(source, sh, pl, turn, prFid, prPr, seFid, sePr, &betrayals)
	// ------ GROUND FIGHT ---------- //
	peace := func(fid1, fid2 int) bool {
		return truces[[2]int{fid1, fid2}] != nil
	}
	betrayedF := func(truster, betrayor int) {
		pt := [2]int{truster, betrayor}
		tr := truces[pt]
		delete(truces, pt)
		tr.DELETE()
		betrayals = append(betrayals, [2]int{betrayor, truster})
	}
	if prFid != 0 && seFid != 0 {
		prPeace := peace(prFid, seFid)
		sePeace := peace(seFid, prFid)
		if shFid != 0 && prPeace && sePeace &&
			peace(shFid, seFid) && peace(seFid, shFid) &&
			peace(prFid, shFid) && peace(shFid, prFid) {

			prPeace = false
			sePeace = false
			for _, pair := range [][2]int{
				[2]int{seFid, shFid}, [2]int{shFid, seFid},
				[2]int{prFid, shFid}, [2]int{shFid, prFid},
				[2]int{prFid, seFid}, [2]int{seFid, prFid},
			} {
				betrayedF(pair[0], pair[1])
			}
			if prPr > 0 {
				prPr -= 1
			}
			if sePr > 0 {
				sePr -= 1
			}
			shSize -= 1
		}
		if !prPeace || !sePeace {
			if prPeace && !sePeace {
				betrayedF(prFid, seFid)
				if prPr > 0 {
					prPr -= 1
				}
			} else if sePeace && !prPeace {
				betrayedF(seFid, prFid)
				if sePr > 0 {
					sePr -= 1
				}
			}
			prPr, sePr = prPr-sePr, sePr-prPr
			if sePr > 0 {
				prPr = sePr
				prFid = seFid
				prPW = sePW
			}
			seFid = 0
			sePr = 0
			sePW = 0
		}
	}
	if shFid != 0 {
		prShPeace := peace(prFid, shFid)
		seShPeace := peace(seFid, shFid)
		shSePeace := peace(shFid, seFid)
		shPrPeace := peace(shFid, prFid)
		if prShPeace && !shPrPeace {
			betrayedF(prFid, shFid)
			if prPr > 0 {
				prPr -= 1
			}
			prShPeace = false
		} else if shPrPeace && !prShPeace {
			betrayedF(shFid, prFid)
			if shSize > 0 {
				shSize -= 1
			}
			shPrPeace = false
		}
		if seShPeace && !shSePeace {
			betrayedF(seFid, shFid)
			if sePr > 0 {
				sePr -= 1
			}
			seShPeace = false
		} else if shSePeace && !seShPeace {
			betrayedF(shFid, seFid)
			if shSize > 0 {
				shSize -= 1
			}
			shSePeace = false
		}
		for shSize > 0 {
			var fightLeft bool
			if prPr > 0 && !shPrPeace {
				fightLeft = true
				prPr -= 1
				shSize -= 1
			}
			if shSize > 0 && sePr > 0 && !shSePeace {
				fightLeft = true
				sePr -= 1
				shSize -= 1
			}
			if !fightLeft {
				break
			}
		}
		if shSize > 0 {
			if !shSePeace && !shPrPeace {
				prFid = shFid
				prPr = shSize
				prPW = sePW
				seFid = 0
				sePr = 0
				sePW = 0
			} else if shSePeace {
				prFid = shFid
				prPr = shSize
				prPW = 0
			} else {
				seFid = shFid
				sePr = shSize
				sePW = 0
			}
		}
		if seFid != 0 && sePr > prPr {
			seFid, prFid = prFid, seFid
			sePr, prPr = prPr, sePr
			sePW, prPW = prPW, sePW
		}
	}
	pl.SetPrimaryFaction(prFid)
	pl.SetPrimaryPresence(prPr)
	pl.SetPrimaryPower(prPW)
	pl.SetSecondaryFaction(prFid)
	pl.SetSecondaryPresence(prPr)
	pl.SetSecondaryPower(sePW)
	return
}

func AllSee(source Source, pl PlanetDat, fid1, fid2, fid3, turn int) {
	pFid, sFid := pl.PrimaryFaction(), pl.SecondaryFaction()
	for _, fid := range []int{fid1, fid2, fid3} {
		if fid == 0 || fid == pFid || fid == sFid {
			continue
		}
		source.UpdatePlanetView(fid, turn, pl)
	}
}

func AllBattleRecords(source Source, lander ShipDat, planet PlanetDat, turn, initPrF, initPrP, initSeF, initSeP int, betrayals *[][2]int) {
	var shFid int
	if lander != nil {
		shFid = lander.FID()
	}
	for _, fid := range []int{initPrF, initSeF, shFid} {
		if fid == 0 {
			continue
		}
		source.NewBattleRecord(lander, fid, turn, initPrF, initPrP, initSeF, initSeP, planet, *betrayals)
	}
}
