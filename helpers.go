package main

import (
	"errors"
	"fmt"
	"html/template"
	"mule/mylog"
	"mule/planetattack/attack"
	"net/http"
	"strconv"
	"strings"
)

const TPDIR = "TEMPLATES/"

var (
	HexPolar = attack.HexPolar
	Log      = mylog.Err
)

func init() {
	mylog.SetErr(DATADIR + "errors.txt")
}

func Apply(t *template.Template, w http.ResponseWriter, d interface{}) {
	err := t.ExecuteTemplate(w, "frame", d)
	if err != nil {
		Log("Error executing", t, "template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MixTemp(fileNames ...string) *template.Template {
	names := make([]string, len(fileNames))
	for i, val := range fileNames {
		names[i] = TPDIR + val + ".html"
	}
	return template.Must(template.New("").Funcs(template.FuncMap{"dict": func(val ...interface{}) (map[string]interface{}, error) {
		if len(val)%2 != 0 {
			err := errors.New("Template dict needs even args")
			Log(err)
			return nil, err
		}
		d := make(map[string]interface{}, len(val)/2)
		for i := 0; i < len(val); i += 2 {
			str, ok := val[i].(string)
			if !ok {
				err := errors.New(fmt.Sprintf("Bad template dict arg", val[i], ": need string keys"))
				Log(err)
				return nil, err
			}
			d[str] = val[i+1]
		}
		return d, nil
	},
	}).ParseFiles(names...))
}

func GetInts(r *http.Request, varNames ...string) (map[string]int, error) {
	m := make(map[string]int, len(varNames))
	badNames := []string{}
	for _, name := range varNames {
		varStr := r.FormValue(name)
		varI, err := strconv.Atoi(varStr)
		if err != nil {
			badNames = append(badNames, fmt.Sprintf("%s:%s", name, err))
		} else {
			m[name] = varI
		}
	}
	if len(badNames) == 0 {
		return m, nil
	}
	return nil, errors.New(fmt.Sprint("GetInts errors:", strings.Join(badNames, "||")))
}

func makeE(v ...interface{}) error {
	return errors.New(fmt.Sprint(v...))
}
