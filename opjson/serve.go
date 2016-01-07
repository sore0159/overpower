package opjson

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	MAXSIZE int64 = 1048576
)

var Log = log.Println

func SetLogger(f func(...interface{})) {
	Log = f
}

func HttpServe(w http.ResponseWriter, j interface{}) (ok bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	sh, ok := MakeShell(j)
	sh.Serve(w)
	return
}

func HttpRead(w http.ResponseWriter, r *http.Request, j interface{}) (ok bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var err error
	defer func() {
		if !ok {
			HttpError(w, 422, "Object read failed:", err)
		}
	}()
	var bytes []byte
	if bytes, err = ioutil.ReadAll(io.LimitReader(r.Body, MAXSIZE)); err != nil {
		Log("JSON R BODY READ ERROR:", err)
		return false
	}
	if err = r.Body.Close(); err != nil {
		Log("JSON R BODY CLOSE ERROR:", err)
		return false
	}
	if err = json.Unmarshal(bytes, j); err != nil {
		Log("JSON UNMARSHAL ERROR FOR:", j, err)
		return false
	}
	return true
}

func HttpError(w http.ResponseWriter, code int, msg ...interface{}) {
	var status string
	switch code {
	case 400: // Status: Zoidberg
		status = "fail"
	case 401: // Status: Gandalf
		status = "fail"
	case 404: // Status: Galt
		status = "error"
	case 500: // Status: Kirk
		status = "error"
	default:
		msg = append([]interface{}{"BAD CODE PASSED TO HTTPERROR:", code, "\nFOR MSG:"}, msg...)
		code = 500
	}
	sh := NewShell()
	sh.Code = code
	sh.Status = status
	sh.Message = fmt.Sprint(msg...)
	sh.Serve(w)
}

func (sh *Shell) Serve(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch sh.Status {
	case "success":
		w.WriteHeader(http.StatusOK)
	case "error", "fail":
		w.WriteHeader(sh.Code)
	default:
		w.WriteHeader(500)
		sh = MakeServerErrShell("BAD SHELL STATUS PASSED TO ServeShell:", sh.Status)
	}
	if err := json.NewEncoder(w).Encode(sh); err != nil {
		Log("JSON SERVESHELL ERROR ENCODE ERROR:", sh, "\n", err)
	}
}
