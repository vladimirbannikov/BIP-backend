package repository

import "errors"

var ErrObjectNotFound = errors.New("not found")

var ErrDuplicateKey = errors.New("duplicate key")
