package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type View struct {
	User  string
	Error string
	path  []string
	App   interface{}
}

func (v *View) SetApp(a interface{}) {
	v.App = a
}

func (v *View) SetError(args ...interface{}) {
	v.Error = fmt.Sprintln(args...)
}

func MakeView(r *http.Request) *View {
	v := View{
		User: CookieUserName(r),
		path: strings.Split(r.URL.Path, "/"),
	}
	return &v
}

func (v *View) Apply(t *template.Template, w http.ResponseWriter) {
	Apply(t, w, v)
}
