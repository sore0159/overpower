package main

import (
	"errors"
	"fmt"
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

/*
func Bail(w http.ResponseWriter, my *mybad.MuleError) bool {
	Log(my)
	http.Error(w, my.Error(), http.StatusInternalServerError)
	return false
}
*/

func (h *Handler) HandleServerError(w http.ResponseWriter, my *mybad.MuleError) {
	Log(my)
	http.Error(w, my.Error(), http.StatusInternalServerError)
}

func (h *Handler) HandleUserError(w http.ResponseWriter, msg string, args ...interface{}) {
	http.Error(w, fmt.Sprintf(msg, args...), http.StatusBadRequest)
}
