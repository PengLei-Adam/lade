package http

import (
	"github.com/PengLei-Adam/lade/app/http/middleware/cors"
	"github.com/PengLei-Adam/lade/app/http/module/demo"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/gin"
	ginSwagger "github.com/PengLei-Adam/lade/framework/middleware/gin-swagger"
	"github.com/PengLei-Adam/lade/framework/middleware/gin-swagger/swaggerFiles"
	"github.com/PengLei-Adam/lade/framework/middleware/static"
)

func Routes(r *gin.Engine) {

	container := r.GetContainer()
	configService := container.MustMake(contract.ConfigKey).(contract.Config)
	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	// 使用cors中间件
	r.Use(cors.Default())

	// 如果配置了swagger, 则显示swagger的中间件
	if configService.GetBool("app.swagger") == true {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	demo.Register(r)
}
