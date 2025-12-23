package loader

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spelens-gud/twig/interfaces/iconfig"
	"gopkg.in/yaml.v3"
)

var _ iconfig.ConfigLoader = (*files)(nil)

// ConfigType 配置文件类型.
type ConfigType uint

const (
	// ConfigTypeUnknown 配置文件类型未知.
	ConfigTypeUnknown ConfigType = iota
	// ConfigTypeJson 配置文件类型JSON.
	ConfigTypeJson
	// ConfigTypeYaml 配置文件类型YAML.
	ConfigTypeYaml
	// ConfigTypeToml 配置文件类型TOML.
	ConfigTypeToml
)

// Unmarshal 反序列化配置数据.
func (t ConfigType) Unmarshal(data []byte, config any) (err error) {
	if !GetEnv().IsProduction() {
		var bf bytes.Buffer
		if err := json.Indent(&bf, data, "", "\t"); err != nil {
			fmt.Printf("\n[CONFIG]:\n%s\n", data)
		} else {
			fmt.Printf("\n[CONFIG]:\n%s\n", bf.String())
		}
	}

	switch t {
	case ConfigTypeYaml:
		return yaml.Unmarshal(data, config)
	case ConfigTypeJson:
		return json.Unmarshal(data, config)
	default:
		return json.Unmarshal(data, config)
	}
}

// files 文件配置加载器.
type files []string

// NewFileLoader function    创建文件配置加载器.
func NewFileLoader(filePaths ...string) iconfig.ConfigLoader {
	return files(filePaths)
}

// Load 加载配置.
func (fs files) Load(configRoot any) (err error) {
	return fs.LoadContext(context.Background(), configRoot)
}

// MustLoad 加载配置,失败则panic.
func (fs files) MustLoad(configRoot any) { must(fs.Load(configRoot)) }

// LoadContext 带上下文加载配置.
func (fs files) LoadContext(ctx context.Context, configRoot any) (err error) {
	return wrapLoad(ctx, configRoot, func(ctx context.Context) (err error) {
		var extDirPath []string
		for _, f := range fs {
			if f = strings.TrimSpace(f); filepath.IsAbs(f) {
				continue
			}
			// 添加调用来源当前目录
			if callerDir, ok := getCaller(); ok {
				nf := filepath.Join(callerDir, f)
				if _, e := os.Stat(nf); e == nil {
					extDirPath = append(extDirPath, nf)
				}
			}

			// 添加二进制当前目录
			nf := filepath.Join(ExecDir(), f)
			if _, e := os.Stat(nf); e == nil {
				extDirPath = append(extDirPath, nf)
			}
		}

		for _, f := range append(fs, extDirPath...) {
			if loadFromFile(f, configRoot) == nil {
				return
			}
		}
		return errors.New("load config from all files failed")
	})
}

func loadFromFile(f string, configRoot any) (err error) {
	var typ ConfigType

	switch {
	case strings.HasSuffix(f, ".json"):
		typ = ConfigTypeJson
	case strings.HasSuffix(f, ".yml"), strings.HasSuffix(f, ".yaml"):
		typ = ConfigTypeYaml
	case strings.HasSuffix(f, ".toml"):
		typ = ConfigTypeToml
	}

	return loadFromBytes(func() ([]byte, error) {
		return os.ReadFile(f)
	}, configRoot, typ)
}
