package cobra

import (
	"log"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/robfig/cron/v3"
)

func (c *Command) SetContainer(container framework.Container) {
	c.container = container
}

// 获取根命令的容器
func (c *Command) GetContainer() framework.Container {
	return c.Root().container
}

// CronSpec 保存Cron命令的信息，用于展示
type CronSpec struct {
	Type        string
	Cmd         *Command
	Spec        string
	ServiceName string
}

func (c *Command) SetParentNull() {
	c.parent = nil
}

// AddCronCommand创建一个cron定时任务，spec为定时描述语句
func (c *Command) AddCronCommand(spec string, cmd *Command) {
	// cron定时器结构初始化保存在根Command内
	root := c.Root()
	if root.Cron == nil {
		root.Cron = cron.New(cron.WithParser(cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)))
		root.CronSpecs = []CronSpec{}
	}
	// 增加说明信息
	root.CronSpecs = append(root.CronSpecs, CronSpec{
		Type: "normal-cron",
		Cmd:  cmd,
		Spec: spec,
	})

	// 制作输入Command的副本
	var cronCmd Command
	ctx := root.Context()
	cronCmd = *cmd
	cronCmd.args = []string{}
	cronCmd.SetParentNull()
	cronCmd.SetContainer(root.GetContainer())
	// 定时器中添加调用函数
	root.Cron.AddFunc(spec, func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println()
			}
		}()
		// 直接运行命令内的函数
		err := cronCmd.ExecuteContext(ctx)
		if err != nil {
			log.Println(err)
		}
	})

}
