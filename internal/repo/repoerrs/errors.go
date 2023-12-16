package repoerrs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrCannotCreate  = errors.New("cannot create")
	ErrCannotGet     = errors.New("cannot get")
	ErrCannotDelete  = errors.New("cannot delete")
)
