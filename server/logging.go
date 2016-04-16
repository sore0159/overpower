package main

import (
	"errors"
	"mule/mybad"
	"mule/mylog"
	"mule/overpower/models"
)

var (
	Check        = mybad.BuildCheck("package", "overpower/server")
	ErrLogger    = mylog.MustErrFile(DATADIR + "errors.txt")
	WarnLogger   = mylog.StockErrorLogger()
	InfoLogger   = mylog.Must(mylog.StockInfoLogger().AddFiles(DATADIR + "info.txt"))
	Ping         = WarnLogger.Ping
	Announce     = InfoLogger.Println
	ErrNoneFound = models.ErrNoneFound
	NewError     = errors.New
)

func Log(err error) {
	if me, ok := err.(*mybad.MuleError); ok {
		WarnLogger.Println(me.MuleError())
		ErrLogger.Println(me.LogError())
	} else {
		WarnLogger.Println(err)
		ErrLogger.Println(err)
	}
}
