package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PengLei-Adam/lade/framework"
)

func FooControllerHandler(ctx *framework.Context) error {
	// 负责通知结束
	finish := make(chan struct{}, 1)
	// 负责通知panic异常
	panicChan := make(chan interface{}, 1)

	durationCtx, cancelFunc := context.WithTimeout(ctx.BaseContext(), time.Second)
	defer cancelFunc()

	go func() {
		defer func() {
			if p := recover(); p != nil {
				panicChan <- p
			}
		}()

		// 业务
		time.Sleep(10 * time.Second)
		ctx.SetOkStatus().Json("ok")

		// 结束的时候通过finish通知父goroutine结束
		finish <- struct{}{}
	}()

	select {
	case p := <-panicChan:
		print(p)
	case <-finish:
		fmt.Println("finish")
	case <-durationCtx.Done():
		ctx.WriterMux().Lock()
		defer ctx.WriterMux().Unlock()
		ctx.SetStatus(500).Json("Time out")
		ctx.SetHasTimeout()

	}

	return nil
}

func UserLoginController(c *framework.Context) error {
	foo, _ := c.QueryString("foo", "def")
	// 等待10s才结束执行
	time.Sleep(10 * time.Second)
	// 输出结果
	c.SetOkStatus().Json("ok, UserLoginController: " + foo)
	return nil
}
