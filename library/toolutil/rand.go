package toolutil

import (
	"math/rand"
	"time"
)

const strList string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStr(length int) string {
	var result []byte

	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, strList[rand.Int63()%int64(len(strList))])
	}
	return string(result)
}

func RandInt(length int) int {
	var code int
	for i := 0; i < 5; i++ {
		code += rand.Intn(10)
	}
	return code
}
