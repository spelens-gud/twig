package kcrypto

import (
	"encoding/base64"
	"encoding/hex"
)

// EncodeKeyHex  将 base64 编码的密钥转换为十六进制.
func EncodeKeyHex(key string) string {
	decoded, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return key
	}
	return hex.EncodeToString(decoded)
}
