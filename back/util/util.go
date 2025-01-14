package util

import (
	"math/rand"
	"time"
)

// 工具函数，若注册时没有传名字，生成随机字符串
func RandomString(n int) string {
	var letters = []byte("fajvhfaufvafbiauAABIFWFIIudfuwfuwdagcqiuetoeh")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// 生成随机数字
func GenerateRandomDigits(length int) string {
	digits := "0123456789"
	result := make([]byte, length)
	rand.Seed(time.Now().Unix())
	for i := 0; i < length; i++ {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}
