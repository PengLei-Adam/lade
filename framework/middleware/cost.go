package middleware

import (
	"log"
	"time"

	"github.com/PengLei-Adam/lade/framework"
)

func Cost() framework.ControllerHandler {
	return func(c *framework.Context) error {
		// 记录开始时间
		start := time.Now()

		// 使用next执行业务逻辑
		c.Next()

		// 记录结束时间
		end := time.Now()
		cost := end.Sub(start)

		log.Printf("api uri %v, cost: %v", c.GetRequest().RequestURI, cost.Seconds())
		return nil
	}
}
