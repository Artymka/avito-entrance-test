package storage

import "errors"

var (
	ErrNoRows     = errors.New("no rows in result")
	ErrUnique     = errors.New("unique constraint violated")
	ErrForeignKey = errors.New("foreign key not found")
)
