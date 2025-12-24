package kv

import (
	"context"
	"strings"
	"time"

	"git.bestfulfill.tech/devops/go-core/interfaces/istore"
)

// option 配置选项.
type option struct {
	keySplit func(string) (string, string)
}

// Option 配置选项.
type Option func(opt *option)

// Options 配置选项列表.
type Options []Option

// init 选项初始化
func (opts Options) init() *option {
	o := &option{
		keySplit: DefaultSplitConfig,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}

// DefaultSplitConfig 默认的keySplit
func DefaultSplitConfig(s string) (string, string) {
	spKey := strings.Split(s, "###")
	if len(spKey) != 2 {
		return s, ""
	}
	return spKey[0], spKey[1]
}

// WithSplitFunc 自定义keySplit
func WithSplitFunc(f func(string) (string, string)) Option {
	return func(opt *option) {
		opt.keySplit = f
	}
}

type storeOption struct {
	ctx                context.Context  // 控制内存回收器生命周期的context
	store              istore.FileStore // 储存器
	writeFileInterval  time.Duration    // 持久化间隔
	clearInterval      time.Duration    // 检查内存驻留超时时间
	defaultExpiration  time.Duration    // 默认缓存时间(缓存时间仅代表多久不从数据源更新数据 但为降级缓存策略 所有数据仍驻留在内存中)
	srcTimeout         time.Duration    // 请求数据源超时时间
	memStoreExpiration time.Duration    // 内存驻留超时
	persistenceFile    string           // 持久化文件名
	blockNums          int              // 分区数
	disableMustGet     bool             // 请求数据源失败且有内存缓存时 不使用降级内存缓存
	lazyUpdate         bool             // 惰性更新 当key过期时不阻塞请求数据源 优先返回旧数据 然后进行异步更新
}
