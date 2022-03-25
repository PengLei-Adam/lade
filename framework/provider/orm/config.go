package orm

import (
	"context"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
	"gorm.io/gorm"
)

// GetBaseConfig 通过configService从配置文件中生成Config
func GetBaseConfig(c framework.Container) *contract.DBConfig {
	configService := c.MustMake(contract.ConfigKey).(contract.Config)
	logService := c.MustMake(contract.LogKey).(contract.Log)
	config := &contract.DBConfig{}
	err := configService.Load("database", config)
	if err != nil {
		logService.Error(context.Background(), "parse database config error", nil)
		return nil
	}
	return config
}

// WithConfigPath 加载配置文件地址
func WithConfigPath(configPath string) contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		// 加载configPath配置路径
		if err := configService.Load(configPath, config); err != nil {
			return err
		}
		return nil
	}
}

// 从gorm.Config中获取配置
func WithGormConfig(gormConfig *gorm.Config) contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		if gormConfig.Logger == nil {
			gormConfig.Logger = config.Logger
		}
		config.Config = gormConfig
		return nil
	}
}

// WithDryRun 设置空跑模式
func WithDryRun() contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		config.DryRun = true
		return nil
	}
}

// WithFullSaveAssociations 设置保存时候关联
func WithFullSaveAssociations() contract.DBOption {
	return func(container framework.Container, config *contract.DBConfig) error {
		config.FullSaveAssociations = true
		return nil
	}
}
