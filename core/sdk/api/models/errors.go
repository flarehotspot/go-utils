package models

import "errors"

var (
	ErrNoValidSession = errors.New("no more available sessions")
)