package command

import "github.com/PengLei-Adam/lade/framework/cobra"

// AddKernelCommands 在跟命令中添加所有命令
func AddKernelCommands(root *cobra.Command) {
	root.AddCommand(DemoCommand)

	//app
	root.AddCommand(initAppCommand())

	// cron
	root.AddCommand(initCronCommand())

	// env
	root.AddCommand(initEnvCommand())

	// config
	root.AddCommand(initAppCommand())
}
