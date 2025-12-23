package loader

import (
	"github.com/spelens-gud/twig/interfaces/iconfig"
	"go.uber.org/atomic"
)

var _ iconfig.EnvChecker = (*Env)(nil)

type Env string

const (
	// EnvUnknown 环境未知.
	EnvUnknown Env = "Unknow"
	// EnvDevelopment 开发环境.
	EnvDevelopment = "Development"
	// EnvTest 测试环境.
	EnvTest = "Test"
	// EnvStaging 预生产环境.
	EnvStaging = "Staging"
	// EnvProduction 生产环境.
	EnvProduction = "Production"
)

// IsDevelopment 是否开发环境.
func (e Env) IsDevelopment() bool { return e == EnvDevelopment || len(e) == 0 }

// IsTesting 是否测试环境.
func (e Env) IsTesting() bool { return e == EnvTest }

// IsPreRelease 是否预发布环境.
func (e Env) IsPreRelease() bool { return e == EnvStaging }

// IsProduction 是否正式环境.
func (e Env) IsProduction() bool { return e == EnvProduction }

// String 获取环境字符串.
func (e Env) String() string {
	switch e {
	case EnvDevelopment:
		return EnvDevelopment
	case EnvTest:
		return EnvTest
	case EnvStaging:
		return EnvStaging
	case EnvProduction:
		return EnvProduction
	default:
		return string(e)
	}
}

var _ iconfig.ConfigSourceLoader = (*configSource)(nil)

type configSource struct {
	atomic.String
	init func() string
}

func NewConfigSource(f func() string) iconfig.ConfigSourceLoader {
	return &configSource{init: f}
}

func (c *configSource) Set(str string) { c.String.Store(str) }

func (c *configSource) Get() string {
	ret := c.Load()
	if len(ret) > 0 {
		return ret
	}
	ret = c.init()
	c.Store(ret)
	return ret
}
