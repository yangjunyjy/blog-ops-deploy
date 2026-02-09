package token

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// randomString 生成随机字符串
func randomString(length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// 如果随机数生成失败，使用时间戳作为后备
		return hex.EncodeToString([]byte(fmt.Sprint(time.Now().UnixNano())))[:length]
	}
	return hex.EncodeToString(bytes)[:length]
}
