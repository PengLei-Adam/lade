package kernel

import (
	"net/http"

	"github.com/PengLei-Adam/lade/framework/gin"
)

// 引擎服务
type LadeKernelService struct {
	engine *gin.Engine
}

// 初始化web引擎服务实例
func NewLadeKernelService(params ...interface{}) (interface{}, error) {
	httpEngine := params[0].(*gin.Engine)
	return &LadeKernelService{engine: httpEngine}, nil
}

// 服务逻辑，返回引擎
func (s *LadeKernelService) HttpEngine() http.Handler {
	return s.engine
}
