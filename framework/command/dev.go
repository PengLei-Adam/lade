package command

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/PengLei-Adam/lade/framework"
	"github.com/PengLei-Adam/lade/framework/cobra"
	"github.com/PengLei-Adam/lade/framework/contract"
	"github.com/PengLei-Adam/lade/framework/util"
	"github.com/fsnotify/fsnotify"
)

// devConfig 代表调试模式的配置信息
type devConfig struct {
	Port    string // 调试模式最终监听的端口，默认为8070
	Backend struct {
		RefreshTime   int    // 文件变更后等待该秒数才更新，默认1s
		Port          string // 后端监听端口，默认8072
		MonitorFolder string // 监听文件夹，默认为AppFolder
	}
	Frontend struct {
		Port string // 前端监听端口，默认8071
	}
}

// 初始化配置文件
func initDevConfig(c framework.Container) *devConfig {
	devConfigPtr := &devConfig{
		Port: "8087",
		Backend: struct {
			RefreshTime   int
			Port          string
			MonitorFolder string
		}{
			1,
			"8072",
			"",
		},
		Frontend: struct{ Port string }{
			"8071",
		},
	}

	configer := c.MustMake(contract.ConfigKey).(contract.Config)
	if configer.IsExist("app.dev.port") {
		devConfigPtr.Port = configer.GetString("app.dev.port")
	}
	if configer.IsExist("app.dev.backend.refresh_time") {
		devConfigPtr.Backend.RefreshTime = configer.GetInt("app.dev.backend.refresh_time")
	}
	monitorFolder := configer.GetString("app.dev.backend.monitor_folder")
	if monitorFolder == "" {
		appService := c.MustMake(contract.AppKey).(contract.App)
		devConfigPtr.Backend.MonitorFolder = appService.AppFolder()
	}

	if configer.IsExist("app.dev.frontend.port") {
		devConfigPtr.Frontend.Port = configer.GetString("app.dev.frontend.port")
	}
	return devConfigPtr
}

// Proxy 代表serve启动的服务器代理
type Proxy struct {
	devConfig   *devConfig   // 配置文件
	proxyServer *http.Server // proxy 服务
	backendPid  int          // 后端服务pid
	frontendPid int          // 前端服务pid
}

// NewProxy 初始化一个Proxy
func NewProxy(c framework.Container) *Proxy {
	devConfig := initDevConfig(c)
	return &Proxy{
		devConfig: devConfig,
	}
}

// newProxyReverseProxy 创建一个反向代理网关
func (p *Proxy) newProxyReverseProxy(frontend, backend *url.URL) *httputil.ReverseProxy {
	if p.frontendPid == 0 && p.backendPid == 0 {
		fmt.Println("前端和后端服务都不存在")
		return nil
	}

	// 后端服务存在
	if p.frontendPid == 0 && p.backendPid != 0 {
		return httputil.NewSingleHostReverseProxy(backend)
	}

	// 前端服务存在
	if p.backendPid == 0 && p.frontendPid != 0 {
		return httputil.NewSingleHostReverseProxy(frontend)
	}

	// 两个都有进程
	director := func(req *http.Request) {
		req.URL.Scheme = backend.Scheme
		req.URL.Host = backend.Host
	}

	// 定义一个NotFoundErr
	NotFoundErr := errors.New("response is 404, need to redirect")
	return &httputil.ReverseProxy{
		Director: director,
		ModifyResponse: func(r *http.Response) error {
			if r.StatusCode == 404 {
				return NotFoundErr
			}
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			// 判断传入的err是否为NotFoundErr，是的话进行前端服务的转发，
			if errors.Is(err, NotFoundErr) {
				httputil.NewSingleHostReverseProxy(frontend).ServeHTTP(w, r)
			}
		},
	}
}

// rebuildBackend 重新编译后端
func (p *Proxy) rebuildBackend() error {
	// 重新编译lade应用
	cmdBuild := exec.Command("./lade", "build", "backend")
	cmdBuild.Stdout = os.Stdout
	cmdBuild.Stderr = os.Stderr
	if err := cmdBuild.Start(); err == nil {
		err = cmdBuild.Wait()
		if err != nil {
			return err
		}
	}
	return nil
}

// restartBackend 启动后端服务
func (p *Proxy) restartBackend() error {
	// 杀死之前的的进程
	if p.backendPid != 0 {
		proc, err := os.FindProcess(p.backendPid)
		if err == nil {
			proc.Kill()
		}
		p.backendPid = 0
	}

	// 设置随机端口，真实后端的端口
	port := p.devConfig.Backend.Port
	ladeAddress := fmt.Sprintf(":" + port)
	// 使用命令行启动后端进程
	cmd := exec.Command("./lade", "app", "start", "--address="+ladeAddress)
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr
	fmt.Println("启动后端服务：", "http://127.0.0.1:"+port)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	p.backendPid = cmd.Process.Pid
	fmt.Println("后端服务pid:", p.backendPid)
	return nil
}

// restartFrontend 启动前端服务
func (p *Proxy) restartFrontend() error {
	// 启动前端调试模式
	// 如果已经开启，有pid存在，则什么都不做
	if p.frontendPid != 0 {
		return nil
	}

	// 开启npm run serve 命令
	port := p.devConfig.Frontend.Port
	path, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	cmd := exec.Command(path, "run", "dev")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s%s", "PORT=", port))
	cmd.Stdout = os.NewFile(0, os.DevNull)
	cmd.Stderr = os.Stderr

	// 使用goroutine启动
	err = cmd.Start()
	fmt.Println("启动前端服务：", "http://127.0.0.1:"+port)
	if err != nil {
		fmt.Println(err)
	}
	p.frontendPid = cmd.Process.Pid
	fmt.Println("前端服务pid:", p.frontendPid)
	return nil
}

// startProxy 根据参数启动前或后端服务
func (p *Proxy) startProxy(startFrontend, startBackend bool) error {
	var backendURL, frontendURL *url.URL
	var err error

	// 启动后端
	if startBackend {
		if err := p.restartBackend(); err != nil {
			return err
		}
	}

	// 启动前端
	if startFrontend {
		if err := p.restartFrontend(); err != nil {
			return err
		}
	}

	// 如果已经启动proxy，返回
	if p.proxyServer != nil {
		return nil
	}

	if frontendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Frontend.Port)); err != nil {
		return err
	}
	if backendURL, err = url.Parse(fmt.Sprintf("%s%s", "http://127.0.0.1:", p.devConfig.Backend.Port)); err != nil {
		return err
	}

	// 设置反向代理
	proxyReverse := p.newProxyReverseProxy(frontendURL, backendURL)
	// 设置服务器
	p.proxyServer = &http.Server{
		Addr:    "127.0.0.1:" + p.devConfig.Port,
		Handler: proxyReverse,
	}

	fmt.Println("代理服务启动：", "http://"+p.proxyServer.Addr)
	// 启动proxy服务
	err = p.proxyServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// monitorBackend 监听应用文件
func (p *Proxy) monitorBackend() error {
	// 监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil
	}
	defer watcher.Close()

	appFolder := p.devConfig.Backend.MonitorFolder
	fmt.Println("监控文件夹：", appFolder)
	filepath.Walk(appFolder, func(path string, info fs.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			return nil
		}
		if util.IsHiddenDirectory(path) {
			return nil
		}
		return watcher.Add(path)
	})

	refreshTime := p.devConfig.Backend.RefreshTime
	t := time.NewTimer(time.Duration(refreshTime) * time.Second)
	t.Stop()
	for {
		select {
		case <-t.C:
			fmt.Println("...检测到文件更新，重启服务开始...")
			if err := p.rebuildBackend(); err != nil {
				fmt.Println("重新编译失败：", err.Error())
			} else {
				if err := p.restartBackend(); err != nil {
					fmt.Println("重新启动失败：", err.Error())
				}
			}
			fmt.Println("...检测到文件更新，重启服务结束...")
			t.Stop()
		case _, ok := <-watcher.Events:
			if !ok {
				continue
			}
			t.Reset(time.Duration(refreshTime) * time.Second)
		case err, ok := <-watcher.Errors:
			if !ok {
				continue
			}
			fmt.Println("监听文件夹错误：", err.Error())
			t.Reset(time.Duration(refreshTime) * time.Second)
		}
	}
	return nil
}

// 初始化Dev命令
func initDevCommand() *cobra.Command {
	devCommand.AddCommand(devBackendCommand)
	devCommand.AddCommand(devFrontendCommand)
	devCommand.AddCommand(devAllCommand)
	return devCommand
}

// devCommand 为调试模式的一级命令
var devCommand = &cobra.Command{
	Use:   "dev",
	Short: "调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		c.Help()
		return nil
	},
}

// devBackendCommand 启动后端调试模式
var devBackendCommand = &cobra.Command{
	Use:   "backend",
	Short: "启动后端调试模式",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		// goroutine中监听后端代码
		go proxy.monitorBackend()
		if err := proxy.startProxy(false, true); err != nil {
			return err
		}
		return nil
	},
}

// devFrontendCommand 启动前端调试模式
var devFrontendCommand = &cobra.Command{
	Use:   "frontend",
	Short: "前端调试模式",
	RunE: func(c *cobra.Command, args []string) error {

		// 启动前端服务
		proxy := NewProxy(c.GetContainer())
		return proxy.startProxy(true, false)

	},
}

var devAllCommand = &cobra.Command{
	Use:   "all",
	Short: "同时启动前端和后端调试",
	RunE: func(c *cobra.Command, args []string) error {
		proxy := NewProxy(c.GetContainer())
		go proxy.monitorBackend()
		if err := proxy.startProxy(true, true); err != nil {
			return err
		}
		return nil
	},
}
