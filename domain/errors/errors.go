package errors

import (
	"github.com/mentoring/toko-backend/domain/errx"
)

type ErrorKind = errx.ErrorKind

const (
	KindNotFound     = errx.KindNotFound
	KindInvalid      = errx.KindInvalid
	KindConflict     = errx.KindConflict
	KindUnauthorized = errx.KindUnauthorized
	KindInternal     = errx.KindInternal
)

type DomainError = errx.DomainError

var (
	New = errx.New
	As  = errx.As
)
