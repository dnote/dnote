package database

import (
	"github.com/pkg/errors"
)

type modelError string

var (
	// ErrNotFound an error that indicates that the given resource is not found
	ErrNotFound error = errors.New("not found")
)
