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

// 认证配置，包括私钥，AccessToken和RefreshToken的过期时间(以秒为单位)
type AuthConfig struct {
	Secret             string `toml:"secret"`
	AccessTokenExpire  int    `toml:"access_token_expire"`
	RefreshTokenLength int    `toml:"refresh_token_length"`
	RefreshTokenExpire int    `toml:"refresh_token_expire"`
}

// 加载配置函数，需要指定路径
func LoadConfig(path string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, fmt.Errorf("读取配置失败: %w", err)
	}
	return &config, nil
}
