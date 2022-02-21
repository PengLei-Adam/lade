package kernel

import (
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/gin"
)

// LadeKernelProvider提供web引擎对象
type LadeKernelProvider struct {
	HttpEngine *gin.Engine
}

// Register注册服务Provider
func (provider *LadeKernelProvider) Register(c framework.Container) framework.NewInstance {
	return NewLadeKernelService
}

// Boot，Web启动前判断是否有外界注入了Engine，用注入的Engine，否则新建Engine
func (provider *LadeKernelProvider) Boot(c framework.Container) error {
	if provider.HttpEngine == nil {
		provider.HttpEngine = gin.Default()
	}
	provider.HttpEngine.SetContainer(c)
	return nil
}

// IsDefer 引擎的初始化我们希望开始就进行初始化
func (provider *LadeKernelProvider) IsDefer() bool {
	return false
}

// Params 参数就是一个HttpEngine
func (provider *LadeKernelProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.HttpEngine}
}

// Name 提供凭证
func (provider *LadeKernelProvider) Name() string {
	return contract.KernelKey
}
