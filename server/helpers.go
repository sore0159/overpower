package main

import (
	"mule/mydb"
	"mule/mylog"
	"mule/myweb"
)

const (
	TPDIR   = "TEMPLATES/"
	DATADIR = "DATA/"
)

var (
	Log     = mylog.QuietErr
	MixTemp = myweb.MakeMixer(TPDIR, map[string]interface{}{
		"link": func(placeholder string) string {
			return "PLACEHOLDER"
		},
		"command": func(placeholder string) string {
			return "PLACEHOLDER"
		},
		"dict": myweb.TemplateDict,
	})
	ValidText = myweb.TextValid
	GetInts   = myweb.GetInts
)

func init() {
	mylog.SetErr(DATADIR + "errors.txt")
	myweb.SetLogger(mylog.QuietErr)
	mydb.SetLogger(mylog.QuietErr)
}
