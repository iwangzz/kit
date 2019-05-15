package main

import "errors"

var (
	ErrInvalidRequestType  = errors.New("invalid request type")
	ErrorBadRequest  = errors.New("invalid request parameter")
)