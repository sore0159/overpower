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
	t2 := t.Funcs(template.FuncMap{
		"link": v.Link,
	})
	Apply(t2, w, v)
}

func (v *View) newpath(n int) string {
	return strings.Join(v.path[:n], "/")
}

func (v *View) Link(str string) string {
	r := []string{""}
	for _, part := range append(v.path, strings.Split(str, "/")...) {
		if part != "" {
			r = append(r, part)
		}
	}
	r = append(r, "")
	return strings.Join(r, "/")
}
