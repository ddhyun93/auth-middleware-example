package utils

import "errors"

var ErrNoTokenGiven = errors.New("no token given")
var ErrUserDeactivated = errors.New("deactivated user")
var ErrExpiredToken = errors.New("token expired")
var ErrInvalidToken = errors.New("invalid token")
