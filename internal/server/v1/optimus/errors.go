package optimus

import "errors"

var (
	ErrOptimusHostNotFound = errors.New("could not find optimus host in shield project metadata")
	ErrOptimusHostInvalid  = errors.New("optimus host is not valid")
)
