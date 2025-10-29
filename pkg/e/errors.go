package e

import "errors"

var (
	ErrInvalidInputFormat = errors.New("invalid input format")
	ErrEmptyRows          = errors.New("empty data for query")
)
