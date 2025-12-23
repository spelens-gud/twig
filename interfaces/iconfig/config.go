package iconfig

import "context"

// ConfigLoader 配置加载器接口.
type ConfigLoader interface {
	// LoadContext 带上下文加载配置.
	LoadContext(ctx context.Context, configRoot any) (err error)
	// Load 加载配置.
	Load(configRoot any) (err error)
	// MustLoad 加载配置,失败则panic.
	MustLoad(configRoot any)
}

// ConfigSourceLoader 配置源加载器接口.
type ConfigSourceLoader interface {
	Set(string)
	Get() string
}
