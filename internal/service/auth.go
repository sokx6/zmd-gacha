package service

import (
	"fmt"

	"zmd-gacha/internal/database"
	"zmd-gacha/internal/utils"
)

func Register(username string, password string, email string) error {
	hashed_pwd, err := utils.HashPWD(password)
	if err != nil {
		return fmt.Errorf("密码哈希失败: %w", err)
	}

	if err = database.RegisterUser(username, hashed_pwd, email); err != nil {
		return fmt.Errorf("数据库存储用户数据失败: %w", err)
	}
	return nil
}
