package errors

import (
	"errors"
)

type DomainError struct {
	ErrType ErrorType
	Err     error
}

type ErrorType int

const (
	InvalidParameter ErrorType = iota
	UnPemitedOperation
	AlreadyExist
	RepositoryError
	QueryError
	QueryDataNotFoundError
	ErrorUnknown
	TransactionError
)

func (e *DomainError) Error() string {
	if e == nil {
		return ""
	}
	return e.Err.Error()
}

func (e *DomainError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

func (e *DomainError) Is(target error) bool {
	if e == nil {
		return false
	}
	var targetErr *DomainError
	if errors.As(target, &targetErr) {
		return e.ErrType == targetErr.ErrType
	}
	return false
}

func (e *DomainError) GetType() ErrorType {
	return e.ErrType
}

func NewDomainError(errType ErrorType, message string) *DomainError {
	return &DomainError{
		ErrType: errType,
		Err:     errors.New(message),
	}
}
