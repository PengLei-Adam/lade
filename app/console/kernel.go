package console

import (
	"github.com/PengLei-Adam/lade/app/console/command/demo"
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/command"
)

// RunCommand is command
//func RunCommand(container framework.Container) error {
//	var rootCmd = &cobra.Command{
//		Use:   "hade",
//		Short: "main",
//		Long:  "hade commands",
//	}
//
//	ctx := commandUtil.RegiestContainer(container, rootCmd)
//
//	hadeCommand.AddKernelCommands(rootCmd)
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
	rootCmd.AddCommand(demo.InitFoo())
}
