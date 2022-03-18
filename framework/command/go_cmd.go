package command

import (
	"log"
	"os"
	"os/exec"

	"github.com/PengLei-Adam/lade/framework/cobra"
)

// go 命令知识运行本地go
var goCommand = &cobra.Command{
	Use:   "go",
	Short: "运行path/go程序，go必须安装",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := exec.LookPath("go")
		if err != nil {
			log.Fatalln("lade go: should install go in your PATH")
		}

		c := exec.Command(path, args...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
		return nil
	},
}
