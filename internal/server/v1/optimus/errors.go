package optimus

import "errors"

var (
	ErrOptimusHostNotFound  = errors.New("optimus hostname doesnot exist")
	ErrOptimusHostNotString = errors.New("optimus hostname doesnot exist")
)
