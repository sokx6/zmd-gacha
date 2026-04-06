package types

import "fmt"

type DatabaseError struct {
	Message string
}

func (e DatabaseError) Error() string {
	return fmt.Sprintf("数据库错误: %s", e.Message)
}

var UserExistsError = DatabaseError{
	Message: "用户已存在",
}

var DatabaseConnectionError = DatabaseError{
	Message: "无法连接到数据库",
}

var DatabaseDefaultError = DatabaseError{
	Message: "数据库操作失败",
}
