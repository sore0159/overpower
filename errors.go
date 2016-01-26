package overpower

import (
	"errors"
	"mule/mybad"
)

var (
	Check        = mybad.BuildCheck("package", "overpower")
	ErrIgnorable = errors.New("something bad happened but we carried on")
	ErrBadArgs   = errors.New("config data for something was bad")
)
