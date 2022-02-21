package contract

import "net/http"

// KernelKey提供kernel服务凭证
const KernelKey = "lade:kernel"

// Kernel接口提供框架最核心的结构
type Kernel interface {
	// 返回net/http启动时需要的Handler接口，实际返回gin.Engine对象
	HttpEngine() http.Handler
}
