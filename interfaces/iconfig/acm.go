package iconfig

import "context"

// AcmConfigLoader 配置加载器接口.
type AcmConfigLoader interface {
	// GetConfig 获取配置.
	GetConfig(context.Context, AcmConfigData) ([]byte, error)
	// GetOriginConfig 获取原始配置.
	GetOriginConfig(context.Context, AcmConfigData) ([]byte, error)
	// UpdateConfig 更新配置.
	UpdateConfig(context.Context, AcmConfigData) error
	// UpdateConfigSet 更新配置集.
	UpdateConfigSet(context.Context, AcmConfigData) error
}

// AcmConfigData 配置数据结构.
type AcmConfigData struct {
	TenantID  string `url:"tenantID" json:"tenantID"`
	DataID    string `url:"dataID" json:"dataID"`
	GroupID   string `url:"groupID" json:"groupID"`
	AppID     string `url:"appID" json:"appID"`         // 部署应用ID
	Namespace string `url:"namespace" json:"namespace"` // 部署命名空间
	Content   string `json:"content" url:"-"`
}
