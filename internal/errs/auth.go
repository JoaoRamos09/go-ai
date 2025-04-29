package errs

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired = errors.New("token expired")
	ErrDecodePassword     = errors.New("error decoding password")
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidTokenClaims = errors.New("invalid token claims")
	ErrInvalidTokenStructure = errors.New("invalid token structure")
	ErrTokenRequired = errors.New("token is required")
	ErrUnauthorized = errors.New("unauthorized")
)
