package loader

import (
	"context"
	"fmt"

	"github.com/spelens-gud/twig/interfaces/iconfig"
)

var _ iconfig.ConfigLoader = (*multi)(nil)

// multi 多配置加载器.
type multi []iconfig.ConfigLoader

// MustLoad 加载配置,失败则panic.
func (m multi) MustLoad(configRoot any) { must(m.Load(configRoot)) }

// Load 加载配置.
func (m multi) Load(configRoot any) (err error) {
	return m.LoadContext(context.Background(), configRoot)
}

// LoadContext 带上下文加载配置.
func (m multi) LoadContext(ctx context.Context, configRoot any) (err error) {
	return wrapLoad(ctx, configRoot, func(ctx context.Context) (err error) {
		for _, loader := range m {
			if err = loader.LoadContext(ctx, configRoot); err == nil {
				return
			}
		}
		fmt.Println("[TWIG] load config from all loader failed")
		return
	})
}

// NewMultiLoader 创建多配置加载器.
func NewMultiLoader(loaders ...iconfig.ConfigLoader) iconfig.ConfigLoader {
	return multi(loaders)
}
