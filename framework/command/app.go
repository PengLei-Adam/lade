package command

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/contract"
)

func initAppCommand() *cobra.Command {
	appCommand.AddCommand(appStartCommand)
	return appCommand
}

// AppCommand 是命令行参数第一级为app的命令，没有实际功能，只是打印帮助文档
var appCommand = &cobra.Command{
	Use:   "app",
	Short: "业务应用控制命令",
	Long:  "业务应用控制命令，包含业务启动，关闭，重启，查询等功能",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Help()
		return nil
	},
}

// appStartCommand启动一个Web服务
var appStartCommand = &cobra.Command{
	Use:   "start",
	Short: "启动一个Web服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取服务容器
		container := cmd.GetContainer()
		// 从服务容器中获取kernel的服务实例
		kernelService := container.MustMake(contract.KernelKey).(contract.Kernel)
		// 获取引擎
		core := kernelService.HttpEngine()

		// 创建一个Server服务
		server := &http.Server{
			Handler: core,
			Addr:    ":8888",
		}

		// 启动服务监听goroutine
		go func() {
			server.ListenAndServe()
		}()

		// 当前的goroutine等待信号量
		quit := make(chan os.Signal)
		// 监控信号：SIGINT, SIGTERM, SIGQUIT
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		// 这里会阻塞当前goroutine等待信号
		<-quit

		// 调用Server.Shutdown graceful结束
		timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(timeoutCtx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}

		return nil
	},
}
