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
	root.AddCommand(initConfigCommand())

	// npm build
	root.AddCommand(initBuildCommand())
	// npm [args...]
	root.AddCommand(npmCommand)
	// go build
	root.AddCommand(goCommand)

	// dev
	root.AddCommand(initDevCommand())

	// cmd
	root.AddCommand(initCmdCommand())

	// provider
	root.AddCommand(initProviderCommand())

	// middleware
	root.AddCommand(initMiddlewareCommand())

	// new
	root.AddCommand(initNewCommand())

	// swagger
	root.AddCommand(initSwaggerCommand())
}
