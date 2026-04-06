package types

import "fmt"

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
