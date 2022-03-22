package command

import (
	"fmt"
	"path/filepath"

	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/swaggo/swag/gen"
)

func initSwaggerCommand() *cobra.Command {
	swaggerCommand.AddCommand(swaggerGenCommand)
	return swaggerCommand
}

var swaggerCommand = &cobra.Command{
	Use:   "swagger",
	Short: "swagger对应命令",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
		}
		return nil
	},
}

// 生成swagger文件
var swaggerGenCommand = &cobra.Command{
	Use:   "gen",
	Short: "生成对应的swagger文件，contain swagger.yaml, doc.go",
	Run: func(cmd *cobra.Command, args []string) {
		container := cmd.GetContainer()
		appService := container.MustMake(contract.AppKey).(contract.App)

		outputDir := filepath.Join(appService.AppFolder(), "http", "swagger")

		conf := &gen.Config{
			// 遍历需要查询注释的目录
			SearchDir: "./app/http/",
			// 不包含哪些文件
			Excludes: "",
			// 输出目录
			OutputDir: outputDir,
			// 整个swagger接口的说明文档注释
			MainAPIFile: "swagger.go",
			// 名字的显示策略，比如首字母大写
			PropNamingStrategy: "",
			// 是否要解析vendor目录
			ParseVender: false,
			// 是否要解析外部依赖库的包
			ParseDependency: false,
			// 是否要解析标准库的包
			ParseInternal: false,
			// 是否要查找markdown文件，这个markdown文件能用来为tag增加说明格式
			MarkdownFilesDir: "",
			// 是否应该在docs.go中生成时间戳
			GeneratedTime: false,
		}
		err := gen.New().Build(conf)
		if err != nil {
			fmt.Println(err)
		}
	},
}
