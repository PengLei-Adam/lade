package http

import (
	"github.com/PengLei-Adam/lade/app/http/module/demo"
	"github.com/PengLei-Adam/lade/framework/gin"
	"github.com/PengLei-Adam/lade/framework/middleware/static"
)

func Routes(r *gin.Engine) {

	r.Use(static.Serve("/", static.LocalFile("./dist", false)))

	demo.Register(r)
}
