package app

import (
	"fmt"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

// LadeAppProvider 提供 App 的具体实现方法
type LadeAppProvider struct {
	BaseFolder string
}

var _ framework.ServiceProvider = &LadeAppProvider{"."}

// Params 获取初始化参数
func (h *LadeAppProvider) Params(container framework.Container) []interface{} {
	return []interface{}{container, h.BaseFolder}
}

// Name方法返回字符串凭证
func (sp *LadeAppProvider) Name() string {
	return contract.AppKey
}

func (sp *LadeAppProvider) Register(c framework.Container) framework.NewInstance {
	return NewLadeApp
}

func (sp *LadeAppProvider) Boot(c framework.Container) error {
	fmt.Println("Lade App service boot")
	return nil
}

func (sp *LadeAppProvider) IsDefer() bool {
	return true
}
