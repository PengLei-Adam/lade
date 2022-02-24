package env

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
	"strings"

	"github.com/PengLei-Adam/lade/framework/contract"
)

// LadeEnv 是 Env 的具体实现
type LadeEnv struct {
	folder string            // .env所在目录
	maps   map[string]string // 保存所有环境变量
}

// NewLadeEnv 有一个参数，.env文件所在的目录
// example: NewLadeEnv("/envfolder/") 会读取文件: /envfolder/.env
// .env的文件格式 FOO_ENV=BAR
func NewLadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewLadeEnv param error")
	}

	folder := params[0].(string)

	// 实例化
	ladeEnv := &LadeEnv{
		folder: folder,
		maps:   map[string]string{"APP_ENV": contract.EnvDevelopment},
	}

	// 解析folder/.env文件
	file := path.Join(folder, ".env")

	// 打开文件.env
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		br := bufio.NewReader(fi)
		for {
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}

			s := bytes.SplitN(line, []byte{'='}, 2)
			if len(s) < 2 {
				continue
			}

			// 保存map
			key := string(s[0])
			val := string(s[1])
			ladeEnv.maps[key] = val
		}
	}

	// 获取当前程序的环境变量，覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) < 2 {
			continue
		}
		ladeEnv.maps[pair[0]] = pair[1]
	}

	return ladeEnv, nil
}

// AppEnv 获取表示当前APP环境的变量APP_ENV
func (en *LadeEnv) AppEnv() string {
	return en.Get("APP_ENV")
}

// IsExist 判断一个环境变量是否有被设置
func (en *LadeEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// Get 获取某个环境变量，如果没有设置，返回""
func (en *LadeEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

// All 获取所有的环境变量，.env和运行环境变量融合后结果
func (en *LadeEnv) All() map[string]string {
	return en.maps
}
