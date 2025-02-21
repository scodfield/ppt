package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// SignSHA256WithKey SHA-256签名
func SignSHA256WithKey(str, key string) string {
	strSign := str + key
	// 计算SHA-256哈希值
	hash := sha256.Sum256([]byte(strSign))
	// 哈希转16进制小写字符串
	signValue := hex.EncodeToString(hash[:])
	return signValue
}
