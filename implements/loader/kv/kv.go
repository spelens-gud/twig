package kv

import (
	"context"

	"github.com/spelens-gud/twig/interfaces/iconfig"
)

var _ iconfig.KVConfigSourceLoader = (*Config)(nil)

// Config 配置.
type Config struct {
	tenantID string                  // 命名空间ID
	cli      iconfig.AcmConfigLoader // 配置加载器
	opt      option                  // 配置加载器选项
}

// NewConfig 创建配置.
func NewConfig(tenantID string, cli iconfig.AcmConfigLoader, opts ...Option) *Config { // 创建配置.
	o := Options(opts).init()
	return &Config{
		tenantID: tenantID,
		cli:      cli,
		opt:      *o,
	}
}

// GetWithContext 获取配置.
func (a *Config) GetWithContext(ctx context.Context, key string) (value string, err error) {
	groupID, dataID := a.opt.keySplit(key)
	ret, err := a.cli.GetOriginConfig(ctx, iconfig.AcmConfigData{
		DataID:   dataID,
		GroupID:  groupID,
		TenantID: a.tenantID,
	})
	if err != nil {
		return
	}
	value = string(ret)
	return
}

// SetWithContext 设置配置.
func (a *Config) SetWithContext(ctx context.Context, key string, value string) (err error) {
	dataID, groupID := a.opt.keySplit(key)
	err = a.cli.UpdateConfig(ctx, iconfig.AcmConfigData{
		DataID:   dataID,
		GroupID:  groupID,
		Content:  value,
		TenantID: a.tenantID,
	})
	return
}
