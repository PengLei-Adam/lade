package app

import (
	"errors"
	"flag"
	"path/filepath"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/util"
	"github.com/google/uuid"
)

type LadeApp struct {
	container  framework.Container
	baseFolder string
	appId      string

	configMap map[string]string // 配置加载
}

func (l LadeApp) AppID() string {
	return l.appId
}

func (l LadeApp) Version() string {
	return "0.0.2"
}

// BaseFolder 表示基础目录，可以代表开发场景的目录，也可以代表运行时候的目录
func (h LadeApp) BaseFolder() string {
	if h.baseFolder != "" {
		return h.baseFolder
	}
	// 如果参数也没有，使用默认的当前路径
	return util.GetExecDirectory()
}

// ConfigFolder  表示配置文件地址
func (h LadeApp) ConfigFolder() string {
	if val, ok := h.configMap["config_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "config")
}

// LogFolder 表示日志存放地址
func (h LadeApp) LogFolder() string {
	if val, ok := h.configMap["log_folder"]; ok {
		return val
	}
	return filepath.Join(h.StorageFolder(), "log")
}

func (h LadeApp) HttpFolder() string {
	if val, ok := h.configMap["http_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "http")
}

func (h LadeApp) ConsoleFolder() string {
	if val, ok := h.configMap["console_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "console")
}

func (h LadeApp) StorageFolder() string {
	if val, ok := h.configMap["storage_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "storage")
}

// ProviderFolder 定义业务自己的服务提供者地址
func (h LadeApp) ProviderFolder() string {
	if val, ok := h.configMap["provider_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "provider")
}

// MiddlewareFolder 定义业务自己定义的中间件
func (h LadeApp) MiddlewareFolder() string {
	if val, ok := h.configMap["middleware_folder"]; ok {
		return val
	}
	return filepath.Join(h.HttpFolder(), "middleware")
}

// CommandFolder 定义业务定义的命令
func (h LadeApp) CommandFolder() string {
	if val, ok := h.configMap["command_folder"]; ok {
		return val
	}
	return filepath.Join(h.ConsoleFolder(), "command")
}

// RuntimeFolder 定义业务的运行中间态信息
func (h LadeApp) RuntimeFolder() string {
	if val, ok := h.configMap["runtime_folder"]; ok {
		return val
	}
	return filepath.Join(h.StorageFolder(), "runtime")
}

// TestFolder 定义测试需要的信息
func (h LadeApp) TestFolder() string {
	if val, ok := h.configMap["test_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "test")
}

func (h LadeApp) AppFolder() string {
	if val, ok := h.configMap["app_folder"]; ok {
		return val
	}
	return filepath.Join(h.BaseFolder(), "app")
}

// NewLadeApp 初始化 LadeApp
func NewLadeApp(params ...interface{}) (interface{}, error) {
	if len(params) != 2 {
		return nil, errors.New("param error")
	}
	// 通过uuid生成appid
	appId := uuid.New().String()

	// 有两个参数，一个是容器，一个是 baseFolder
	container := params[0].(framework.Container)
	baseFolder := params[1].(string)
	// 如果没有设置，则使用参数
	if baseFolder == "" {
		flag.StringVar(&baseFolder, "base_folder", "", "base_folder参数, 默认为当前路径")
		flag.Parse()
	}
	return &LadeApp{baseFolder: baseFolder, container: container, appId: appId}, nil
}

func (app LadeApp) LoadAppConfig(kv map[string]string) {
	for k, v := range kv {
		app.configMap[k] = v
	}
}
