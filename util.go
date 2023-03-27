package aliyundrive_open

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// randomString 生成随机字符串
func randomString(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

// 合并自定义的字符串类型
func joinCustomString[T fmt.Stringer](items []T, separator string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return fmt.Sprint(items[0])
	default:
		var b strings.Builder
		b.WriteString(fmt.Sprint(items[0]))
		for _, s := range items[1:] {
			b.WriteString(separator)
			b.WriteString(fmt.Sprint(s))
		}
		return b.String()
	}
}
