package token

import "errors"

var ErrExpiredToken = errors.New("token is expired")
var ErrInvalidToken = errors.New("token is invalid")
