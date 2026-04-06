package database

import (
	"fmt"
	"zmd-gacha/internal/config"
	"zmd-gacha/internal/models"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	DB   *gorm.DB
	Cfg  *config.DataBaseConfig
	UIDs map[uint]bool
}

var defaultDB *Database

func NweDatabase(Cfg config.DataBaseConfig) *Database {
	return &Database{
		Cfg:  &Cfg,
		UIDs: make(map[uint]bool),
	}
}

func Init(cfg config.DataBaseConfig) error {
	db := NweDatabase(cfg)
	if err := db.InitDB(); err != nil {
		return err
	}
	defaultDB = db
	return nil
}

func Get() (*Database, error) {
	if defaultDB == nil || defaultDB.DB == nil {
		return nil, fmt.Errorf("数据库未初始化")
	}
	return defaultDB, nil
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

	database.DB.AutoMigrate(&models.User{}, &models.Character{}, &models.RefreshToken{})

	var uids []uint
	database.DB.Model(&models.User{}).Pluck("uid", &uids)
	for _, uid := range uids {
		database.UIDs[uid] = true
	}

	return nil
}
