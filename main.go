package main

import (
	"github.com/PengLei-Adam/lade/app/console"
	"github.com/PengLei-Adam/lade/app/http"
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/provider/app"
	"github.com/PengLei-Adam/lade/framework/provider/config"
	"github.com/PengLei-Adam/lade/framework/provider/distributed"
	"github.com/PengLei-Adam/lade/framework/provider/env"
	"github.com/PengLei-Adam/lade/framework/provider/id"
	"github.com/PengLei-Adam/lade/framework/provider/log"
	"github.com/PengLei-Adam/lade/framework/provider/orm"
	"github.com/PengLei-Adam/lade/framework/provider/trace"
	"github.com/PengLei-Adam/lade/provider/kernel"
)

func main() {
	// 初始化服务容器
	container := framework.NewLadeContainer()
	// 绑定App服务提供者
	container.Bind(&app.LadeAppProvider{})
	// 绑定分布式服务
	container.Bind(&distributed.LocalDistributedProvider{})
	// 绑定环境变量服务
	// 未设置folder？？
	container.Bind(&env.LadeEnvProvider{})
	// 绑定加载配置文件的服务
	container.Bind(&config.LadeConfigProvider{})
	container.Bind(&id.LadeIDProvider{})
	container.Bind(&trace.LadeTraceProvider{})
	container.Bind(&log.LadeLogServiceProvider{})
	container.Bind(&orm.GormProvider{})
	// 后续绑定其他服务提供者

	// 将HTTP引擎初始化，并且作为服务提供者绑定到服务容器中
	if engine, err := http.NewHttpEngine(container); err == nil {
		container.Bind(&kernel.LadeKernelProvider{HttpEngine: engine})
	}

	// 运行root命令
	console.RunCommand(container)

}
