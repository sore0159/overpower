package main

import (
	"mule/overpower"
	"time"
)

func AutoTimer() {
	for {
		now := time.Now()
		then := now.Round(time.Hour)
		if hour := then.Hour(); hour > 22 {
			then = then.Add(time.Hour * time.Duration(23+(24-hour)))
		} else {
			then = then.Add(time.Hour * time.Duration(23-hour))
		}
		then = then.Add(time.Minute * time.Duration(2))
		InfoLog("Starting autotimer:", now, "\n SLEEP TILL:", then)
		dur := then.Sub(now)
		time.Sleep(dur)
		now = time.Now()
		if now.Hour() != 23 {
			Log("AUTO RUN ERROR: SLEEP DID NOT REACH HOUR 23")
			continue
		} else {
			InfoLog("Autotimer woke:", now)
		}
		DBLOCK = true
		time.Sleep(5 * time.Minute)
		games, ok := OPDB.AllGames()
		if !ok {
			Log("AUTO RUN ERROR FETCHING GAMES")
			continue
		}
		var count int
		countChan := make(chan byte)
		wkDay := int(now.Weekday())
		for _, g := range games {
			if g.Turn() < 1 {
				continue
			}
			days := g.AutoDays()
			if days[wkDay] {
				var run bool
				if free := g.FreeAutos(); free > 0 {
					g.SetFreeAutos(free - 1)
				} else {
					run = true
				}
				count++
				go func(g overpower.Game, run bool, done chan byte) {
					if run {
						InfoLog("AUTO RUNNING GAME", g.Gid())
						if !OPDB.AutoRunGameTurn(g) {
							Log("ERROR AUTO RUNNING GAME", g.Gid())
						}
					} else {
						if !OPDB.UpdateGame(g) {
							Log("ERROR UPDATING FREEAUTOS FOR GAME", g.Gid())
						}
					}
					done <- 0
				}(g, run, countChan)
			}
		}
		for count > 0 {
			<-countChan
			count -= 1
		}
		DBLOCK = false
	}
}
