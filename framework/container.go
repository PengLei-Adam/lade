package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，相同凭证对应提供者替换为新的提供者，返回 error
	Bind(provider ServiceProvider) error
	// 判断关键字凭证是否已经绑定提供者
	IsBind(key string) bool
	// 根据关键字凭证获取一个服务
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会 panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的 params 参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// LadeContainer 是服务容器的具体实现
type LadeContainer struct {
	Container // 强制要求 ladeContainer 实现 Container 接口
	// providers 存储注册的服务提供者，key 为字符串凭证
	providers map[string]ServiceProvider
	// instance 存储具体的实例，key 为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex
}

func NewLadeContainer() *LadeContainer {
	return &LadeContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (lade *LadeContainer) PrintProviders() []string {
	ret := []string{}
	for _, provider := range lade.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind将服务容器和关键字做了绑定
func (lade *LadeContainer) Bind(provider ServiceProvider) error {
	// 涉及写操作，加锁
	lade.lock.Lock()
	defer lade.lock.Unlock()

	key := provider.Name()

	lade.providers[key] = provider

	// if provider is not defer
	if !provider.IsDefer() {
		// 预处理
		if err := provider.Boot(lade); err != nil {
			return err
		}

		// 实例化方法
		params := provider.Params(lade)
		method := provider.Register(lade)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		lade.instances[key] = instance
	}
	return nil
}

// 调用内部make方法
func (lade *LadeContainer) Make(key string) (interface{}, error) {
	return lade.make(key, nil, false)
}

func (lade *LadeContainer) MustMake(key string) interface{} {
	serv, err := lade.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

// MakeNew方式使用内部make初始化
func (lade *LadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return lade.make(key, params, true)
}

// 真正实例化一个服务
func (lade *LadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	lade.lock.RLock()
	defer lade.lock.RUnlock()
	// 查询是否已经注册了服务的提供者，每个注册，则返回错误
	sp := lade.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return lade.newInstance(sp, params)
	}

	// 不需要强制重新实例化，取容器中已有的实例
	if ins, ok := lade.instances[key]; ok {
		return ins, nil
	}
	// 容器中还未实例化，进行一次实例化
	inst, err := lade.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	lade.instances[key] = inst
	return inst, nil
}

func (lade *LadeContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(lade); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(lade)
	}
	method := sp.Register(lade)
	inst, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return inst, err
}

func (lade *LadeContainer) findServiceProvider(key string) ServiceProvider {
	lade.lock.RLock()
	defer lade.lock.RUnlock()
	sp, ok := lade.providers[key]
	if !ok {
		return nil
	}
	return sp
}
