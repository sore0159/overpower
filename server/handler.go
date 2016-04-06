package main

import (
	"fmt"
	"html/template"
	sq "mule/mydb/sql"
	"mule/myweb"
	"mule/overpower"
	"mule/overpower/models"
	"mule/users"
	"net/http"
)

type Handler struct {
	TitleBar    bool
	User        users.User
	LoggedIn    bool
	Error       string
	CommandData [2]int
	M           *models.Manager
	*myweb.Handler
}

func MakeHandler(w http.ResponseWriter, r *http.Request) *Handler {
	var user users.User
	test, ok, err := USERREG.IsLoggedIn(w, r)
	if my, bad := Check(err, "handler creation problem", "path", r.URL.Path); bad {
		Log(my)
	} else if ok {
		user = test
	}
	return &Handler{TitleBar: true, User: user, LoggedIn: ok, M: OPDB.NewManager(), Handler: myweb.MakeHandler(r)}
}

func (h *Handler) Apply(t *template.Template, w http.ResponseWriter) {
	err := h.Handler.Apply(w, t.Funcs(template.FuncMap{
		"link":    h.Link,
		"command": h.Command,
	}), "frame", h)
	if my, bad := Check(err, "handler apply failure", "path", h.Path, "template", t); bad {
		Log(my)
	}
}

func (h *Handler) SetCommand(g overpower.GameDat) {
	h.CommandData[0] = g.GID()
	h.CommandData[1] = g.Turn()
}

func (h *Handler) Command(str string) string {
	return fmt.Sprintf("/overpower/command/%d/%d/%s", h.CommandData[0], h.CommandData[1], str)
}

func (h *Handler) SetError(args ...interface{}) {
	h.Error = fmt.Sprintln(args...)
}

func (h *Handler) GID(gid int) sq.Condition {
	return sq.EQ("gid", gid)
}
func (h *Handler) FID(gid, fid int) sq.Condition {
	return sq.AND(sq.EQ("gid", gid), sq.EQ("fid", fid))
}
