package attack

import (
	"fmt"
	"mule/mylog"
)

var Log = mylog.Err

func LogF(v ...interface{}) {
	Log(v)
	panic(fmt.Sprintln(v...))
}
