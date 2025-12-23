package kgo

import "unsafe"

func SetStringIfEmpty(in *string, backs ...string) {
	out := *in
	if len(out) > 0 {
		return
	}
	*in = GetFirstValidString(backs...)
}

// GetFirstValidString 初始化字符串
func GetFirstValidString(strings ...string) string {
	for _, str := range strings {
		if len(str) > 0 {
			return str
		}
	}
	return ""
}

// UnsafeBytes2string 字节转字符串
func UnsafeBytes2string(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
