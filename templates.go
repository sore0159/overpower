package main

import (
	"html/template"
	"net/http"
)

const TPDIR = "TEMPLATES/"

var TPLOGIN = MixTemp("frame", "login")

func LoginTP(w http.ResponseWriter, d interface{}) {
	err := TPLOGIN.ExecuteTemplate(w, "frame", d)
	if err != nil {
		Log("Error executing Login template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MixTemp(fileNames ...string) *template.Template {
	names := make([]string, len(fileNames))
	for i, val := range fileNames {
		names[i] = TPDIR + val + ".html"
	}
	return template.Must(template.ParseFiles(names...))
}

type TPContext struct {
	template string
	w        http.ResponseWriter
}

func NewTPContext() *TPContext {
	return &TPContext{
	//
	}
}

func MakeContext(w http.ResponseWriter, tpName string) *TPContext {
	d := NewTPContext()
	d.template = tpName
	d.w = w
	return d
}

func (c *TPContext) Exe() {
	var TP *template.Template
	switch c.template {
	case "login":
		TP = TPLOGIN
	default:
		Log("TPContext Exe with bad template:", c.template)
		http.Error(c.w, "Templating Error!", http.StatusInternalServerError)
		return
	}
	err := TP.ExecuteTemplate(c.w, "frame", c)
	if err != nil {
		Log("Error executing template", c.template, ":", err)
		http.Error(c.w, err.Error(), http.StatusInternalServerError)
	}
}
