package demo

import (
	"time"

	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/gin"
	"github.com/PengLei-Adam/lade/framework/provider/redis"
)

// DemoRedis redis路由handler方法
func (api *DemoApi) DemoRedis(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)

	// 初始化一个redis
	redisService := c.MustMake(contract.RedisKey).(contract.RedisService)
	client, err := redisService.GetClient(redis.WithConfigPath("cache.default"), redis.WithRedisConfig(func(options *contract.RedisConfig) {
		options.MaxRetries = 3
	}))
	if err != nil {
		logger.Error(c, err.Error(), nil)
		c.AbortWithError(50001, err)
	}
	// 向redis写入一个值{foo: bar},有效期1小时
	if err := client.Set(c, "foo", "bar", 1*time.Hour).Err(); err != nil {
		c.AbortWithError(500, err)
		return
	}
	// 获取foo对应的值
	val := client.Get(c, "foo").String()
	logger.Info(c, "redis get", map[string]interface{}{
		"val": val,
	})

	// 删除foo
	if err := client.Del(c, "foo").Err(); err != nil {
		c.AbortWithError(500, err)
		return
	}

	c.JSON(200, "ok")
}

// DemoCache cache的简单例子
func (api *DemoApi) DemoCache(c *gin.Context) {
	logger := c.MustMakeLog()
	logger.Info(c, "request start", nil)
	// 初始化cache服务
	cacheService := c.MustMake(contract.CacheKey).(contract.CacheService)
	// 设置key为foo
	err := cacheService.Set(c, "foo", "bar", 1*time.Hour)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	// 获取key为foo
	val, err := cacheService.Get(c, "foo")
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	logger.Info(c, "cache get", map[string]interface{}{
		"val": val,
	})
	// 删除key为foo
	if err := cacheService.Del(c, "foo"); err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.JSON(200, "ok")
}
