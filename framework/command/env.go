package command

import (
	"fmt"

	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/util"
)

func initEnvCommand() *cobra.Command {
	envCommand.AddCommand(envListCommand)
	return envCommand
}

// 获取当前的App环境的命令
var envCommand = &cobra.Command{
	Use:   "env",
	Short: "获取当前App环境",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)
		// 打印环境
		fmt.Println("environment:", envService.AppEnv())
	},
}

// envListCommand 获取所有App环境变量
var envListCommand = &cobra.Command{
	Use:   "list",
	Short: "获取所有的环境变量",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		envService := container.MustMake(contract.EnvKey).(contract.Env)

		envs := envService.All()
		outs := [][]string{}
		for k, v := range envs {
			outs = append(outs, []string{k, v})
		}
		util.PrettyPrint(outs)
	},
}
