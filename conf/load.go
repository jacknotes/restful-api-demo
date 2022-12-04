package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/caarlos0/env/v6"
)

// LoadConfigFromToml 从toml中添加配置文件, 并初始化全局对象
func LoadConfigFromToml(filePath string) error {
	// 新建一个空的Config对象
	cfg := NewDefaultConfig()
	// toml.DecodeFile解码文件是否符合toml格式，并且将值赋值给cfg对象
	if _, err := toml.DecodeFile(filePath, cfg); err != nil {
		return err
	}
	// 加载全局配置单例
	SetGlobalConfig(cfg)
	return nil
}

// LoadConfigFromEnv 从环境变量中加载配置
func LoadConfigFromEnv() error {
	// 新建一个空的Config对象
	cfg := NewDefaultConfig()
	// env.Parse解析变量并且赋值给cfg对象
	if err := env.Parse(cfg); err != nil {
		return err
	}
	// 加载全局配置单例
	SetGlobalConfig(cfg)
	return nil
}
