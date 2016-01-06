package opjson

import (
	"encoding/json"
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
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		Log("JSON MARSHAL ERROR:", err)
		return false
	}
	return true
}

func HttpRead(w http.ResponseWriter, r *http.Request, j interface{}) (ok bool) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var err error
	defer func() {
		if !ok {
			HttpError(w, err, 422)
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

func HttpError(w http.ResponseWriter, err error, code int) (ok bool) {
	w.WriteHeader(code)
	if err = json.NewEncoder(w).Encode(err); err != nil {
		Log("JSON ERROR ENCODE ERROR:", err)
		return false
	}
	return true
}
