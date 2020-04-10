package store

import "errors"

var (
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	// ErrIllegalArgs ...
	ErrIllegalArgs = errors.New("illegal arguments exception")
)
