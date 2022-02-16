package demo

import (
	"fmt"

	"github.com/PengLei-Adam/lade/framework"
)

type DemoService struct {
	Service

	c framework.Container
}

// 初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	// 参数展开，第一个参数转换为Container
	c := params[0].(framework.Container)

	fmt.Println("new demo service")

	return &DemoService{c: c}, nil
}

// 实现接口,服务具体业务逻辑
func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}
