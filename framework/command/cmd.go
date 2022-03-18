package command

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/util"
	"github.com/jianfengye/collection"
	"github.com/pkg/errors"
)

// 底层级别的控制台命令，列出命令，通过模板创建新命令
func initCmdCommand() *cobra.Command {
	cmdCommand.AddCommand(cmdListCommand)
	cmdCommand.AddCommand(cmdCreateCommand)
	return cmdCommand
}

// 控制台命令
var cmdCommand = &cobra.Command{
	Use:   "command",
	Short: "控制台命令相关",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

// cmdListCommand 列出所有控制台命令
var cmdListCommand = &cobra.Command{
	Use:   "list",
	Short: "列出所有控制台命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		cmds := cmd.Root().Commands()
		ps := [][]string{}
		for _, c := range cmds {
			line := []string{c.Name(), c.Short}
			ps = append(ps, line)
		}
		util.PrettyPrint(ps)
		return nil
	},
}

// cmdCreateCommand 创建一个业务控制台命令
var cmdCreateCommand = &cobra.Command{
	Use:     "new",
	Aliases: []string{"create", "init"}, // 设置别名为create init
	Short:   "创建一个控制台命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		fmt.Println("开始创建控制台命令")
		var name string
		var folder string
		{
			prompt := &survey.Input{
				Message: "请输入控制台命令名称:",
			}
			err := survey.AskOne(prompt, &name)
			if err != nil {
				return err
			}
		}
		{
			prompt := &survey.Input{
				Message: "请输入文件夹名称(默认: 同控制台命令):",
			}
			err := survey.AskOne(prompt, &folder)
			if err != nil {
				return err
			}
		}
		if folder == "" {
			folder = name
		}

		// 判断文件不存在
		app := container.MustMake(contract.AppKey).(contract.App)

		pFolder := app.CommandFolder()
		subFolders, err := util.SubDir(pFolder)
		if err != nil {
			return err
		}
		subColl := collection.NewStrCollection(subFolders)
		if subColl.Contains(folder) {
			fmt.Println("目录名称已经存在")
			return nil
		}

		// 开始创建文件
		if err := os.Mkdir(filepath.Join(pFolder, folder), 0700); err != nil {
			return err
		}

		// 创建title这个模板方法
		funcs := template.FuncMap{"title": strings.Title}
		{
			// 创建name.go
			file := filepath.Join(pFolder, folder, name+".go")
			f, err := os.Create(file)
			if err != nil {
				return errors.Cause(err)
			}

			// 使用contractTmp模板初始化template，并且上这个模板支持title方法
			t := template.Must(template.New("cmd").Funcs(funcs).Parse(cmdTmpl))
			// 将name传递进入到template中渲染，并输出到contract.go中
			if err := t.Execute(f, name); err != nil {
				return errors.Cause(err)
			}
		}

		fmt.Println("创建新命令行工具成功，路径:", filepath.Join(pFolder, folder))
		fmt.Println("请记得开发完成后将命令行工具挂载到 console/kernel.go")
		return nil
	},
}

// 命令行模板
var cmdTmpl string = `package {{.}}

improt (
	"fmt"

	"github.com/PengLei-Adam/lade/framework/cobra"
)

var {{.|title}}Command = &cobra.Command{
	Use: "{{.}}",
	Short: "{{.}}",
	RunE: func(c *cobra.Command, args []string) error {
        container := c.GetContainer()
		fmt.Println(container)
		return nil
	},
}
`
