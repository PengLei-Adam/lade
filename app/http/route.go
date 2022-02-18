package http

import (
	"github.com/PengLei-Adam/lade/app/http/module/demo"
	"github.com/PengLei-Adam/lade/framework/gin"
)

func Routes(r *gin.Engine) {
	r.Static("/dist/", "./dist/")

	demo.Register(r)
}
