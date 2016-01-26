package main

import (
	"mule/myweb"
)

const (
	TPDIR = "TEMPLATES/"
)

var (
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
	GetIntsIf = myweb.GetIntsIf
)
