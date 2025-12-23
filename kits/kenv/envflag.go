// Package kenv 提供环境变量相关的工具函数.
// 包含环境变量标志判断、dotenv 加载等功能.
package kenv

import (
	"os"
	"sync"
)

// BoolOnceFrom 创建一次性布尔环境变量判断函数.
func BoolOnceFrom(key string) func() bool {
	var (
		o sync.Once
		b *bool
	)
	return func() bool {
		o.Do(func() {
			tmp := IsNotEmpty(key)
			b = &tmp
		})
		return *b
	}
}

// BoolFrom 创建布尔环境变量判断函数.
func BoolFrom(key string) func() bool {
	return func() bool {
		return IsNotEmpty(key)
	}
}

// IsNotEmpty 判断环境变量是否非空.
func IsNotEmpty(key string) bool {
	return !IsEmpty(key)
}

// IsEmpty 判断环境变量是否为空.
func IsEmpty(key string) bool {
	return len(os.Getenv(key)) == 0
}

// Get 获取环境变量值.
func Get(key string) string {
	return os.Getenv(key)
}

// GetDefault 获取环境变量值,为空时返回默认值.
func GetDefault(key, defaultValue string) string {
	if v := os.Getenv(key); len(v) > 0 {
		return v
	}
	return defaultValue
}

// Set 设置环境变量.
func Set(key, value string) error {
	return os.Setenv(key, value)
}

// Unset 删除环境变量.
func Unset(key string) error {
	return os.Unsetenv(key)
}
