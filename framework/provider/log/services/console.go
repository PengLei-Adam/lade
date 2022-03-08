package services

import (
	"os"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

// 控制台输出日志
type LadeConsoleLog struct {
	LadeLog
}

func NewLadeConsoleLog(params ...interface{}) (interface{}, error) {
	c := params[0].(framework.Container)
	level := params[1].(contract.LogLevel)
	ctxFielder := params[2].(contract.CtxFielder)
	formatter := params[3].(contract.Formatter)

	log := &LadeConsoleLog{}

	log.SetLevel(level)
	log.SetCtxFielder(ctxFielder)
	log.SetFormatter(formatter)

	// 设置输出端为控制台
	log.SetOutput(os.Stdout)
	log.c = c
	return log, nil
}
