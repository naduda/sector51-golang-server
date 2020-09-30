package store

import "errors"

var (
	// ErrReject ...
	ErrReject = errors.New("operation was rejected")
	// ErrRecordNotFound ...
	ErrRecordNotFound = errors.New("record not found")
	// ErrIllegalArgs ...
	ErrIllegalArgs = errors.New("illegal arguments exception")
)
