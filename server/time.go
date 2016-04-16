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
		Announce("Starting autotimer:", now, "\n SLEEP TILL:", then)
		dur := then.Sub(now)
		time.Sleep(dur)
		now = time.Now()
		if now.Hour() != 23 {
			ErrLogger.Println("AUTO RUN ERROR: SLEEP DID NOT REACH HOUR 23")
			continue
		} else {
			Announce("Autotimer woke:", now)
		}
		DBLOCK = true
		time.Sleep(5 * time.Minute)
		m := OPDB.NewManager()
		games, err := m.Game().Select()
		if my, bad := Check(err, "resource failure in autotimer", "resource", "games"); bad {
			Log(my)
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
				if free := g.FreeAutos(); free > 0 {
					g.SetFreeAutos(free - 1)
				} else {
					count++
					go func(g overpower.GameDat, done chan byte) {
						Announce("AUTO RUNNING GAME", g.GID())
						logE, failE := OPDB.SourceTransact(g.GID(), overpower.RunGameTurn)
						if my, bad := Check(failE, "failure on auto-running turn", "gid", g.GID()); bad {
							Log(my)
						}
						if logE != nil {
							Log(logE)
						}
						done <- 0
					}(g, countChan)
				}
			}
		}
		err = m.Close()
		if my, bad := Check(err, "autorun game freeturn inc failure"); bad {
			Log(my)
		}
		for count > 0 {
			<-countChan
			count -= 1
		}
		DBLOCK = false
	}
}
