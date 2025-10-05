package utils

import "errors"

var (
	// ErrDivisionByZero はゼロ除算エラー
	ErrDivisionByZero = errors.New("division by zero")
)
