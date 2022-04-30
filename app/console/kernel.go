package console

import (
	"github.com/PengLei-Adam/lade/app/console/command/foo"
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/command"
)

// RunCommand is command
//func RunCommand(container framework.Container) error {
//	var rootCmd = &cobra.Command{
//		Use:   "lade",
//		Short: "main",
//		Long:  "lade commands",
//	}
//
//	ctx := commandUtil.RegiestContainer(container, rootCmd)
//
//	ladeCommand.AddKernelCommands(rootCmd)
//
//	// rootCmd.AddCronCommand("* * * * *", command.DemoCommand)
//
//	return rootCmd.ExecuteContext(ctx)
//}

func RunCommand(container framework.Container) error {
	// 根Command
	var rootCommand = &cobra.Command{
		Use:   "lade",
		Short: "lade 命令",
		Long:  "lade 框架提供的命令行工具，执行框架自带命令，也能方便编写业务命令",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不显示cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	// 根命令设置container
	rootCommand.SetContainer(container)
	// 绑定框架提供的所有命令
	command.AddKernelCommands(rootCommand)
	// 绑定业务自定义的命令，在consle/comamnd/XXX/XXX.go中
	AddAppCommand(rootCommand)

	//执行RootCommand
	return rootCommand.Execute()

}

// 绑定业务的命令
func AddAppCommand(rootCmd *cobra.Command) {
	//  demo 例子
	rootCmd.AddCommand(foo.InitFoo())

	// 每秒调用一次Foo命令
	//rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)

	// 启动一个分布式任务调度，调度的服务名称为init_func_for_test，每个节点每5s调用一次Foo命令，抢占到了调度任务的节点将抢占锁持续挂载2s才释放
	//rootCmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}
