package optimus

import "github.com/goto/dex/pkg/errors"

var (
	ErrOptimusHostNotFound = errors.New("No Optimus jobs found in this project")
	ErrOptimusHostInvalid  = errors.New("Optimus host is not valid")
)
