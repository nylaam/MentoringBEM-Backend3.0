package errx

import "errors"

type ErrorKind int

const (
	KindNotFound ErrorKind = iota + 1
	KindInvalid
	KindConflict
	KindUnauthorized
	KindInternal
)

type DomainError struct {
	Kind    ErrorKind
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func New(kind ErrorKind, code, message string) *DomainError {
	return &DomainError{
		Kind:    kind,
		Code:    code,
		Message: message,
	}
}

// As mengekstrak DomainError dari error chain (Slide 13 di index.md)
func As(err error) (*DomainError, bool) {
	var de *DomainError
	if errors.As(err, &de) {
		return de, true
	}
	return nil, false
}
