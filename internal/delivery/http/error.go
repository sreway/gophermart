package http

import (
	"errors"
)

var (
	ErrSignedJWT  = errors.New("failed signed jwt")
	ErrDecodeBody = errors.New("failed decode body")
	ErrEncodeData = errors.New("failed encode data")
)
