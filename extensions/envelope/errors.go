package envelope

import "errors"

var (
	ErrSignatureNotFound = errors.New(`signature not found`)

	ErrSignatureCheckFailed = errors.New(`check signature failed`)

	ErrDeadlineExpired = errors.New(`deadline expired`)

	ErrCheckSignatureFailed = errors.New(`check signature failed`)

	ErrTxAlreadyExecuted = errors.New(`tx already executed`)

	ErrInvalidMethod = errors.New(`invalid method in envelope`)

	ErrInvalidChannel = errors.New(`invalid channel in envelope`)
)
