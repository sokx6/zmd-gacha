package database

import (
	"fmt"
	"zmd-gacha/internal/config"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB  *gorm.DB
	Cfg *config.DataBaseConfig
}

func NweDatabase(Cfg config.DataBaseConfig) *Database {
	return &Database{
		Cfg: &Cfg,
	}
}

func (database *Database) InitDB() error {
	var dialector gorm.Dialector
	var err error
	switch database.Cfg.Driver {
	case "sqlite":
		dialector = sqlite.Open(database.Cfg.DSN)
	case "postgres", "postgresql":
		dialector = postgres.Open(database.Cfg.DSN)
	case "mysql":
		dialector = mysql.Open(database.Cfg.DSN)
	default:
		return fmt.Errorf("未知的数据库类型")
	}

	database.DB, err = gorm.Open(dialector)
	if err != nil {
		return fmt.Errorf("数据库连接错误: %w", err)
	}

	return nil
}
