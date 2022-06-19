package server

import "errors"

var (
	// Key not found
	ErrNoKey = errors.New("no key found")
	// Keys expired
	ErrKeyExpired = errors.New("key expired")
	// Value is not an integer
	ErrValueIsNotInteger = errors.New("value is not an integer")
)
