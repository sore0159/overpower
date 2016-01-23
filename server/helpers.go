package main

import (
	"mule/mybad"
	"mule/mydb"
	"mule/mylog"
	"mule/myweb"
)

const (
	TPDIR   = "TEMPLATES/"
	DATADIR = "DATA/"
)

var (
	Logger     = mylog.MustErrFile(DATADIR + "errors.txt")
	InfoLogger = mylog.Must(mylog.StockInfoLogger().AddFiles(DATADIR + "info.txt"))
	Log        = Logger.Println
	InfoLog    = InfoLogger.Println
	MixTemp    = myweb.MakeMixer(TPDIR, map[string]interface{}{
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
	GetIntsIf = myweb.GetIntsIf
	Check     = mybad.BuildCheck("package", "overpower/server")
)

func init() {
	mydb.SetLogger(Log)
}
