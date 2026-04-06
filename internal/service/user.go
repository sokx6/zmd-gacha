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

	db, err := database.Get()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %w", err)
	}

	if err = db.RegisterUser(username, hashed_pwd, email); err != nil {
		return fmt.Errorf("数据库存储用户数据失败: %w", err)
	}
	return nil
}
