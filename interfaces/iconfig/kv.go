package iconfig

import (
	"context"
	"time"
)

// KVConfigLoader 键值配置加载器接口.
type KVConfigLoader interface {
	// Get 获取配置项.
	Get(string) (string, error)
	// Set 设置配置项.
	Set(string, string, time.Duration) error
	// GetWithContext 带上下文获取配置项.
	GetWithContext(context.Context, string) (string, error)
	// SetWithContext 带上下文设置配置项.
	SetWithContext(context.Context, string, string, time.Duration) error
}

// KVConfigSourceLoader 键值配置源加载器接口.
type KVConfigSourceLoader interface {
	// GetWithContext 带上下文获取配置项.
	GetWithContext(context.Context, string) (string, error)
	// SetWithContext 带上下文设置配置项.
	SetWithContext(context.Context, string, string) error
}
