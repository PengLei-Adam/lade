package config

import (
	"path/filepath"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

type LadeConfigProvider struct{}

// Register registe a new function for make a service instance
func (provider *LadeConfigProvider) Register(c framework.Container) framework.NewInstance {
	return NewLadeConfig
}

// Boot will called when the service instantiate
func (provider *LadeConfigProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *LadeConfigProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *LadeConfigProvider) Params(c framework.Container) []interface{} {
	appService := c.MustMake(contract.AppKey).(contract.App)
	envService := c.MustMake(contract.EnvKey).(contract.Env)
	env := envService.AppEnv()
	// 配置文件夹地址
	configFolder := appService.ConfigFolder()
	envFolder := filepath.Join(configFolder, env)
	return []interface{}{c, envFolder, envService.All()}
}

/// Name define the name for this service
func (provider *LadeConfigProvider) Name() string {
	return contract.ConfigKey
}
