package http

import (
	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/gin"
)

func NewHttpEngine(container framework.Container) (*gin.Engine, error) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.SetContainer(container)

	r.Use(gin.Recovery())

	Routes(r)
	return r, nil
}
