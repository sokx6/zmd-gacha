package utils

import "math/rand"

// 生成随机9位UID
func GenerateUID() uint {
	// 生成一个随机数，范围在100000000到999999999之间
	return 100000000 + uint(rand.Intn(900000000))
}
