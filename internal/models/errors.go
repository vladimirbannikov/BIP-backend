package models

import "errors"

var ErrNotFound = errors.New("item not found")

var ErrInvalidInput = errors.New("invalid input")

var ErrConflict = errors.New("item already exists")

var ErrBadCredentials = errors.New("bad password or login")

var ErrInvalidToken = errors.New("invalid token")

var ErrTokenExpired = errors.New("token expired")

var ErrNo2FAOrExpired = errors.New("2FA code is not present or expired")

var Err2FACodeInvalid = errors.New("client 2FA code is invalid")

var ErrInvalidQuestionOrAnswer = errors.New("invalid question or answer")
