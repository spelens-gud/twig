package istore

import "time"

// Store 存储器接口.
type Store interface {
	// Get 获取缓存数据.
	Get(string) (string, error)
	// Set 设置缓存数据.
	Set(string, string) error
	// Add 添加缓存数据.
	Add(string, string) error
	// Delete 删除缓存数据.
	Delete(string) error
}

// FileStore 文件存储器接口.
type FileStore interface {
	// Delete 删除缓存数据.
	Delete(string)
	// Set 添加缓存数据.
	Set(string, any, time.Duration)
	// Get 获取缓存数据.
	Get(string) (any, bool)
	// LoadFile 加载文件
	LoadFile(string) error
	// SaveFile 保存文件
	SaveFile(string) error
}
