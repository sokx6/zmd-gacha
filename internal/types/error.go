package types

import (
	"errors"
	"fmt"
)

type UserError struct {
	Message string
}

func (e UserError) Error() string {
	return fmt.Sprintf("用户错误: %s", e.Message)
}

type DatabaseError struct {
	Message string
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("数据库错误: %s", e.Message)
}

var UserExistsError = UserError{
	Message: "用户已存在",
}

var UserNotFoundError = UserError{
	Message: "用户未找到",
}

var PasswordIncorrectError = UserError{
	Message: "密码错误",
}

var DatabaseGetError = DatabaseError{
	Message: "获取数据库实例失败",
}

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
