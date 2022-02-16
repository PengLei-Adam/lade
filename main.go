package main

import (
	"context"
	"coredemo/framework"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	core := framework.NewCore()
	registerRouter(core)

	server := &http.Server{
		Addr:    ":8080",
		Handler: core,
	}

	quit := make(chan os.Signal)

	// 所有连接处理完成，再关闭
	// 运行监听
	go func() {
		server.ListenAndServe()
	}()

	// 等待退出信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	// 退出
	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal("Server shutdown: ", err)
	}
}
