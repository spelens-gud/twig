package iconfig

// EnvChecker 环境变量检查器.
type EnvChecker interface {
	// IsDevelopment 是否开发环境.
	IsDevelopment() bool
	// IsTesting 是否测试环境.
	IsTesting() bool
	// IsPreRelease 是否预发布环境.
	IsPreRelease() bool
	// IsProduction 是否正式环境.
	IsProduction() bool
	// String 获取环境字符串.
	String() string
}
