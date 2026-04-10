package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// 配置文件模型
type Config struct {
	App      AppConfig
	Database DataBaseConfig
	Auth     AuthConfig
}

// 应用配置，只是监听的地址和端口
type AppConfig struct {
	Host string `toml:"host"`
	Port uint16 `toml:"port"`
}

// 数据库配置
type DataBaseConfig struct {
	Driver string `toml:"driver"`
	DSN    string `toml:"dsn"`
}

type AuthConfig struct {
	Secret string
}

// 加载配置函数，需要指定路径
func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}
	return &config, nil
}
