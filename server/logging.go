package main

import (
	"mule/mybad"
	"mule/mylog"
	"mule/overpower/models"
	"net/http"
)

var (
	Check        = mybad.BuildCheck("package", "overpower/server")
	ErrLogger    = mylog.MustErrFile(DATADIR + "errors.txt")
	WarnLogger   = mylog.StockErrorLogger()
	InfoLogger   = mylog.Must(mylog.StockInfoLogger().AddFiles(DATADIR + "info.txt"))
	Ping         = WarnLogger.Ping
	Announce     = InfoLogger.Println
	ErrNoneFound = models.ErrNoneFound
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

func Bail(w http.ResponseWriter, my *mybad.MuleError) bool {
	Log(my)
	http.Error(w, my.Error(), http.StatusInternalServerError)
	return false
}
