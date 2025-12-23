// Package loader 提供配置加载器的具体实现.
// 支持文件加载、多源加载等功能.
package loader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spelens-gud/twig/interfaces/iconfig"
)

var (
	pkgPath = func() string {
		_, f, _, _ := runtime.Caller(0)
		return filepath.Dir(f)
	}()

	configCachePath = filepath.Join(ExecDir(), "config-cache", ExecName()+".config.json")
)

const (
	defaultConfigTemplateName = "config.template.json"
	defaultLocalConfigName    = "config.local.json"
)

// NewConfigLoaderFromEnv 根据环境创建配置加载器.
func NewConfigLoaderFromEnv(extLoaders ...iconfig.ConfigLoader) iconfig.ConfigLoader {
	loaders := make(multi, 0)

	if GetEnv().IsDevelopment() {
		loaders = append(loaders, NewFileLoader(defaultLocalConfigName))
	}

	return NewMultiLoader(append(loaders, extLoaders...)...)
}

// getCaller 获取调用者信息.
func getCaller() (dir string, ok bool) {
	for i := 1; i < 10; i++ {
		_, callFile, _, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.HasPrefix(callFile, pkgPath) || strings.Contains(callFile, "/sync/") {
			continue
		}
		return filepath.Dir(callFile), true
	}
	return
}

// wrapLoad 封装加载配置方法.
func wrapLoad(ctx context.Context, configRoot any, load func(ctx context.Context) (err error)) (err error) {
	if GetEnv().IsDevelopment() {
		if dir, ok := getCaller(); ok {
			f := filepath.Join(dir, defaultConfigTemplateName)
			if e := writeConfig(configRoot, f); e != nil {
				fmt.Printf("[TWIG] write config template error: %v\n", e)
			}
		}
	}

	if GetEnv().IsProduction() {
		defer func() {
			if err == nil {
				if e := writeConfig(configRoot, configCachePath); e != nil {
					fmt.Printf("[TWIG] write config cache error: %v\n", e)
				}
				return
			}

			fmt.Printf("[TWIG] load config error: %v, try fallback cache\n", err)
			if err = NewFileLoader(configCachePath).LoadContext(ctx, configRoot); err != nil {
				fmt.Printf("[TWIG] fallback cache error: %v\n", err)
			}
		}()
	}

	if err = load(ctx); err != nil {
		return
	}
	return
}

// writeConfig 写入配置文件.
func writeConfig(configRoot any, file string) (err error) {
	data, err := json.MarshalIndent(configRoot, "", "\t")
	if err != nil {
		return
	}
	if err = os.MkdirAll(filepath.Dir(file), 0775); err != nil {
		return
	}
	return os.WriteFile(file, data, 0664)
}

// loadFromBytes 从字节加载配置.
func loadFromBytes(f func() ([]byte, error), configRoot any, typ ConfigType) (err error) {
	configData, err := f()
	if err != nil {
		fmt.Printf("[TWIG] load config error: %v\n", err)
		return
	}

	if err = typ.Unmarshal(configData, configRoot); err != nil {
		fmt.Printf("[TWIG] unmarshal config error: %v\n", err)
		return
	}
	fmt.Println("[TWIG] config loaded")
	return
}

// must 保证函数调用成功.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
