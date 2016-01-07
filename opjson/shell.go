package opjson

import (
	"encoding/json"
	"fmt"
)

// Shell for encoding per the JSEND api spec
// http://labs.omniti.com/labs/jsend
type Shell struct {
	Status  string          `json:"status"`
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func NewShell() *Shell {
	return &Shell{}
}

// Status: Solo
// We're fine.  We're all fine here now, thank you.  How are you?
func MakeShell(obj interface{}) (*Shell, bool) {
	raw, err := json.Marshal(obj)
	if err != nil {
		Log("ERROR JSON MARSHALLING FOR SHELL:", obj, err)
		return MakeServerErrShell("ERROR MARSHALLING OBJECT", obj, "TO JSON:", err), false
	}
	return &Shell{
		Status: "success",
		Data:   raw,
	}, true
}

// Status: Galt
// Who is John Galt?
func Make404Shell(msg ...interface{}) *Shell {
	return &Shell{
		Status:  "error",
		Message: fmt.Sprint(msg...),
		Code:    404,
	}
}

// Status: Gandalf
// YOU SHALL NOT PASS!
func MakeNoAuthShell(msg ...interface{}) *Shell {
	return &Shell{
		Status:  "fail",
		Message: fmt.Sprint(msg...),
		Code:    401,
	}
}

// Status: Zoidberg
// Your request was bad and you should feel bad!
func MakeBadReqShell(msg ...interface{}) *Shell {
	return &Shell{
		Status:  "fail",
		Message: fmt.Sprint(msg...),
		Code:    400,
	}
}

// Status: Kirk
// WHAT IS LOVE?!
func MakeServerErrShell(msg ...interface{}) *Shell {
	return &Shell{
		Status:  "error",
		Message: fmt.Sprint(msg...),
		Code:    500,
	}
}
