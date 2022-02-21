package demo

import (
	"log"

	"github.com/PengLei-Adam/lade/framework/cobra"
)

// InitFoo 初始化Foo命令
func InitFoo() *cobra.Command {
	FooCommand.AddCommand(Foo1Command)
	return FooCommand
}

// FooCommand 代表Foo命令
var FooCommand = &cobra.Command{
	Use:     "foo",
	Short:   "foo的简要说明",
	Long:    "foo的长说明",
	Aliases: []string{"fo", "f"},
	Example: "foo命令的例子",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		log.Println(container)
		return nil
	},
}

// Foo1Command 代表Foo命令的子命令Foo1
var Foo1Command = &cobra.Command{
	Use:     "foo1",
	Short:   "foo1的简要说明",
	Long:    "foo1的上说明",
	Aliases: []string{"fo1", "f1"},
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		log.Println(container)
		return nil
	},
}
