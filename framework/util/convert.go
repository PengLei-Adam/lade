package util

import (
	"unsafe"

	"github.com/PengLei-Adam/lade/framework/gin"
	ggin "github.com/gin-gonic/gin"
)

func ConvertGinHanderFunc(f ggin.HandlerFunc) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c := (*ggin.Context)(unsafe.Pointer(ctx.GinBaseContext))
		f(c)
	}
}
