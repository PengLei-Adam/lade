package cobra

import "github.com/PengLei-Adam/lade/framework"

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// 获取根命令的容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}
