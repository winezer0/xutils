package hashutils

import (
	"crypto/md5"
	"encoding/hex"
)

// GetStrHashShort 获取字符串的短 MD5 哈希值（前8位）
func GetStrHashShort(s string) string {
	return GetStrHash(s)[:8]
}

// GetStrHash 获取字符串的哈希值
func GetStrHash(s string) string {
	hash := md5.Sum([]byte(s))
	return hex.EncodeToString(hash[:])
}
