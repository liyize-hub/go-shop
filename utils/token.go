package utils

import (
	"crypto/rand"
	"fmt"
)

// 生成8位的字符串
func RandToken() string {
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
