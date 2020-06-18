package models

import (
	"errors"
)

var (
	ErrNotExists     = errors.New("doesn't exists")
	ErrConflict      = errors.New("new data conflicts with old data")
)
