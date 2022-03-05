package command

import (
	"fmt"

	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/kr/pretty"
)

// initConfigCommand 获取配置相关的命令
func initConfigCommand() *cobra.Command {
	configCommand.AddCommand(configGetCommand)
	return configCommand
}

// configCommand 获取配置相关信息
var configCommand = &cobra.Command{
	Use:   "config",
	Short: "获取配置相关信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

// configGetCommand 获取某个配置信息
var configGetCommand = &cobra.Command{
	Use:   "get",
	Short: "获取某个配置信息",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		configService := container.MustMake(contract.ConfigKey).(contract.Config)
		if len(args) != 1 {
			fmt.Println("参数错误")
			return nil
		}
		configPath := args[0]
		val := configService.Get(configPath)
		if val == nil {
			fmt.Println("配置路径 ", configPath, " 不存在")
			return nil
		}

		fmt.Printf("%# v\n", pretty.Formatter(val))
		return nil
	},
}
