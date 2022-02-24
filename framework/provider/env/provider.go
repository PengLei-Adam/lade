package env

import (
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

type LadeEnvProvider struct {
	Folder string
}

// Register registe a new function for make a service instance
func (provider *LadeEnvProvider) Register(c framework.Container) framework.NewInstance {
	return NewLadeEnv
}

// Boot will called when the service instantiate
func (provider *LadeEnvProvider) Boot(c framework.Container) error {
	app := c.MustMake(contract.AppKey).(contract.App)
	provider.Folder = app.BaseFolder()
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *LadeEnvProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *LadeEnvProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.Folder}
}

/// Name define the name for this service
func (provider *LadeEnvProvider) Name() string {
	return contract.EnvKey
}
