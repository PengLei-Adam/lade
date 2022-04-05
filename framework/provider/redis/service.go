package redis

import (
	"sync"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/go-redis/redis/v8"
)

// LadeRedis redis服务的结构实现
type LadeRedis struct {
	container framework.Container
	clients   map[string]*redis.Client // key为uniqKey，缓存每个连接

	lock *sync.RWMutex
}

// NewLadeRedis 实例化Client
func NewLadeRedis(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)

	clients := make(map[string]*redis.Client)
	lock := &sync.RWMutex{}

	return &LadeRedis{
		container: container,
		clients:   clients,
		lock:      lock,
	}, nil
}

// GetClient 获取Client实例
func (app *LadeRedis) GetClient(option ...contract.RedisOption) (*redis.Client, error) {
	// 读取默认配置
	config := GetBaseConfig(app.container)

	// option对默认配置进行修改
	for _, opt := range option {
		if err := opt(app.container, config); err != nil {
			return nil, err
		}
	}

	// 如果最终的config没有设置dsn，则生成dsn
	key := config.UniqKey()

	// 判断是否已有实例化的redis.Client
	app.lock.RLock()
	if db, ok := app.clients[key]; ok {
		app.lock.RUnlock()
		return db, nil
	}
	app.lock.RUnlock()

	// 没有实例化redis.Client，则实例化
	app.lock.Lock()
	defer app.lock.Unlock()

	// 实例化client
	client := redis.NewClient(config.Options)

	// 挂载到map中
	app.clients[key] = client

	return client, nil

}
