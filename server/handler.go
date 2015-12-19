package main

import (
	"fmt"
	"html/template"
	"mule/myweb"
	"mule/overpower"
	"mule/users"
	"net/http"
)

type Handler struct {
	TitleBar    bool
	User        users.User
	LoggedIn    bool
	Error       string
	CommandData [2]int
	*myweb.Handler
}

func MakeHandler(w http.ResponseWriter, r *http.Request) *Handler {
	var user users.User
	test, ok := USERREG.IsLoggedIn(w, r)
	if ok {
		user = test
	}
	return &Handler{TitleBar: true, User: user, LoggedIn: ok, Handler: myweb.MakeHandler(r)}
}

func (h *Handler) Apply(t *template.Template, w http.ResponseWriter) {
	h.Handler.Apply(w, t.Funcs(template.FuncMap{
		"link":    h.Link,
		"command": h.Command,
	}), "frame", h)
}

func (h *Handler) SetCommand(g overpower.Game) {
	h.CommandData[0] = g.Gid()
	h.CommandData[1] = g.Turn()
}

func (h *Handler) Command(str string) string {
	return fmt.Sprintf("/overpower/command/%d/%d/%s", h.CommandData[0], h.CommandData[1], str)
}

func (h *Handler) SetError(args ...interface{}) {
	h.Error = fmt.Sprintln(args...)
}
