package errors

import (
	"errors"
	"fmt"
)

// ErrorCode 错误码类型
type ErrorCode int

// 错误码定义
const (
	ErrCodeInvalidIP ErrorCode = 1000 + iota
	ErrCodeDatabaseError
	ErrCodeCacheError
	ErrCodeInternalError
	ErrCodeInvalidRequest
)

// AppError 应用错误
type AppError struct {
	Code    ErrorCode
	Message string
	Err     error
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code: %d, message: %s, error: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// Unwrap 实现错误包装接口
func (e *AppError) Unwrap() error {
	return e.Err
}

// New 创建新的错误
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// NewWithError 创建带有底层错误的错误
func NewWithError(code ErrorCode, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Is 检查错误类型
func Is(err error, code ErrorCode) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code == code
	}
	return false
}

// GetCode 获取错误码
func GetCode(err error) ErrorCode {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return ErrCodeInternalError
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return err.Error()
}

// 预定义错误
var (
	ErrInvalidIP      = New(ErrCodeInvalidIP, "无效的IP地址")
	ErrDatabaseError  = New(ErrCodeDatabaseError, "数据库错误")
	ErrCacheError     = New(ErrCodeCacheError, "缓存错误")
	ErrInternalError  = New(ErrCodeInternalError, "内部错误")
	ErrInvalidRequest = New(ErrCodeInvalidRequest, "无效的请求")
)
