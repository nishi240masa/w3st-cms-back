package errors

import (
	"errors"
	"fmt"
)

type DomainError struct {
	ErrType ErrorType
	Message string
	Err     error
}

type ErrorType int

const (
	InvalidParameter ErrorType = iota
	UnPermittedOperation
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
	if e.Message != "" {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
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
	if e == nil || target == nil {
		return false
	}

	var t *DomainError
	ok := errors.As(target, &t)
	if !ok {
		return false
	}

	return e.ErrType == t.ErrType
}

func (e *DomainError) GetType() ErrorType {
	return e.ErrType
}

// 最初に作るとき用
func NewDomainError(errType ErrorType, err error) *DomainError {
	if err == nil {
		err = errors.New("unknown error")
	}
	return &DomainError{
		ErrType: errType,
		Err:     err,
	}
}

// メッセージも付けたいとき用
func NewDomainErrorWithMessage(errType ErrorType, message string) *DomainError {
	return &DomainError{
		ErrType: errType,
		Message: message,
		Err:     errors.New(message),
	}
}

// 既存エラーに追加情報を付けたいとき用
func WrapDomainError(message string, err error) *DomainError {
	if err == nil {
		return nil
	}
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		// すでにDomainErrorなら、上からメッセージだけ追加して積む
		return &DomainError{
			ErrType: domainErr.ErrType,
			Message: message,
			Err:     domainErr,
		}
	}
	// 普通のエラーなら、Unknown扱いでラップ
	return &DomainError{
		ErrType: ErrorUnknown,
		Message: message,
		Err:     err,
	}
}
