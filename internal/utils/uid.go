package utils

import "math/rand"

// 生成随机9位UID
func GenerateUID() uint {
	// 生成一个随机数，范围在1000000001到9999999999之间
	return 1000000001 + uint(rand.Intn(8999999998))
}
