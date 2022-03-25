package trace

import (
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

type LadeTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *LadeTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewLadeTraceService
}

// Boot will called when the service instantiate
func (provider *LadeTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *LadeTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *LadeTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

/// Name define the name for this service
func (provider *LadeTraceProvider) Name() string {
	return contract.TraceKey
}
