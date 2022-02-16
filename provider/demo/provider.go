package demo

import (
	"fmt"

	"github.com/PengLei-Adam/lade/framework"
)

//服务提供者示例
type DemoServiceProvider struct {
}

var _ framework.ServiceProvider = &DemoServiceProvider{}

// Name方法返回字符串凭证
func (sp *DemoServiceProvider) Name() string {
	return Key
}

func (sp *DemoServiceProvider) Register(c framework.Container) framework.NewInstance {
	return NewDemoService
}

func (sp *DemoServiceProvider) Boot(c framework.Container) error {
	fmt.Println("demo service boot")
	return nil
}

func (sp *DemoServiceProvider) IsDefer() bool {
	return true
}

func (sp *DemoServiceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
