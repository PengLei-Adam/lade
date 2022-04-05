package contract

import (
	"fmt"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/go-redis/redis/v8"
)

const RedisKey = "lade:redis"

// RedisOption 代表初始化时候的选项
type RedisOption func(container framework.Container, config *RedisConfig) error

// RedisService 表示一个redis服务
type RedisService interface {
	// GetClient 获取redis连接实例
	GetClient(option ...RedisOption) (*redis.Client, error)
}

// RedisConfig 封装redis.Options配置结构
type RedisConfig struct {
	*redis.Options
}

// UniqKey 唯一标识一个RedisConfig配置
func (config *RedisConfig) UniqKey() string {
	return fmt.Sprintf("%v_%v_%v_%v", config.Addr, config.DB, config.Username, config.Network)
}
