package envelop

import "errors"

var (
	ErrSignatureNotFound = errors.New(`signature not found`)

	ErrSignatureCheckFailed = errors.New(`check signature failed`)

	ErrDecodeEnvelopFailed = errors.New(`decod envelope failed`)

	ErrDeadlineExpired = errors.New(`deadline expired`)
)
