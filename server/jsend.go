package main

import (
	"mule/jsend"
	"mule/mybad"
	sq "mule/mydb/sql"
	"net/http"
)

// /overpower/json/...
func apiJSON(w http.ResponseWriter, r *http.Request) {
	h := MakeHandler(w, r)
	switch r.Method {
	case "GET":
		h.ApiJSONget(w, r)
	case "PUT":
		h.ApiJSONput(w, r)
	default:
		JSONUserError(w, "UNSUPPORTED JSON METHOD", KV{"method", r.Method})
	}
}

// Type KV is a convienience struct for vardic key-val pair args
type KV struct {
	Key   string
	Value interface{}
}

func SQLAND(args ...KV) sq.Condition {
	sqArgs := make([]sq.P, len(args))
	for i, arg := range args {
		sqArgs[i] = sq.P{arg.Key, arg.Value}
	}
	return sq.AllEQ(sqArgs...)

}

func JSONSuccess(w http.ResponseWriter, obj interface{}) {
	err := jsend.Success(w, obj)
	if my, bad := Check(err, "jsend success failure", "obj", obj); bad {
		Log(my)
	}
}

func JSONUserError(w http.ResponseWriter, msg string, args ...KV) {
	data := make(map[string]interface{}, len(args)+1)
	for _, kv := range args {
		data[kv.Key] = kv.Value
	}
	data["message"] = msg
	err := jsend.Fail(w, 400, data)
	if my, bad := Check(err, "jsend fail failure", "data", data); bad {
		Log(my)
	}
}

func JSONServerError(w http.ResponseWriter, my *mybad.MuleError) {
	Log(my)
	if my, bad := Check(jsend.Kirk(w), "api kirk failure"); bad {
		Log(my)
	}
}
