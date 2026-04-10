package types

import (
	"errors"
	"fmt"
)

var InvaildTokenError = errors.New("令牌无效")

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"err"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("应用错误: code=%d, message=%s, err=%v", e.Code, e.Message, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}
