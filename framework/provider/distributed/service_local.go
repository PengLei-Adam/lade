package distributed

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/contract"
)

// LocalDistributedService 代表lade框架的App实现
type LocalDistributedService struct {
	container framework.Container
}

// NewLocalDistributedService 初始化本地分布式服务
func NewLocalDistributedService(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("params error")
	}

	// 输入必须是一个参数container
	container := params[0].(framework.Container)
	return &LocalDistributedService{container: container}, nil
}

// Select 分布式选择器
func (s LocalDistributedService) Select(serviceName string, appID string, holdTime time.Duration) (selectAppID string, err error) {
	appService := s.container.MustMake(contract.AppKey).(contract.App)
	runtimeFolder := appService.RuntimeFolder()
	// 设置欲锁文件的路径，包含服务名称
	lockFile := filepath.Join(runtimeFolder, "distribute_"+serviceName)
	// 打开文件
	lock, err := os.OpenFile(lockFile, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}

	// 尝试独占文件锁
	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	// 没有抢到文件锁
	if err != nil {
		// 读取被选择的appid
		selectAppIDByt, err := ioutil.ReadAll(lock)
		if err != nil {
			return "", err
		}
		return string(selectAppIDByt), err
	}

	// 在一段时间内，选择有效，其他节点不能再抢占
	go func() {
		defer func() {
			// 释放文件锁
			syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)
			// 释放文件
			lock.Close()
			// 删除文件
			os.Remove(lockFile)
		}()
		// 创建选择有效的计时器
		timer := time.NewTimer(holdTime)
		// 等待计时器结束
		<-timer.C

	}()
	// 抢占到了，将抢占到的appID写入文件
	if _, err := lock.WriteString(appID); err != nil {
		return "", err
	}
	return appID, nil
}
