package id

import (
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

type LadeIDProvider struct {
}

// Register registe a new function for make a service instance
func (provider *LadeIDProvider) Register(c framework.Container) framework.NewInstance {
	return NewLadeIDService
}

// Boot will called when the service instantiate
func (provider *LadeIDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *LadeIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *LadeIDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this service
func (provider *LadeIDProvider) Name() string {
	return contract.IDKey
}
