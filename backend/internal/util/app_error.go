package util

import (
	"errors"
	"fmt"
)

// AppError 统一应用错误：携带业务错误码与上下文信息。
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d message=%s cause=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d message=%s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error { return e.Err }

// NewAppError 创建应用错误。
func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// WrapAppError 包装底层错误并透传。
func WrapAppError(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// AsAppError 提取 AppError，若不是则包装为内部错误。
func AsAppError(err error) *AppError {
	var ae *AppError
	if errors.As(err, &ae) {
		return ae
	}
	return &AppError{Code: 50000, Message: "internal server error", Err: err}
}

// 仓储层哨兵错误。
var (
	ErrNotFound       = errors.New("record not found")
	ErrDuplicate      = errors.New("duplicate record")
	ErrConflict       = errors.New("state conflict")
	ErrNoRowsAffected = errors.New("no rows affected")
)
