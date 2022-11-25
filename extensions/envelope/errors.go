package envelope

import "errors"

var (
	ErrSignatureNotFound = errors.New(`signature not found`)

	ErrSignatureCheckFailed = errors.New(`check signature failed`)

	ErrDeadlineExpired = errors.New(`deadline expired`)

	ErrCheckSignatureFailed = errors.New(`check signature failed`)
)
