package optimus

import "errors"

var (
	ErrOptimusHostNotFound  = errors.New("could not find optimus host in shield project metadata")
	ErrOptimusHostNotString = errors.New("optimus host is not a valid string")
)
