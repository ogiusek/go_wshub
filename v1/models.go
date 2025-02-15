package wshub

import (
	"errors"
)

var (
	ErrIdIsTaken   error = errors.New("this id is already taken")
	ErrIdIsMissing error = errors.New("this id is missing")
)
