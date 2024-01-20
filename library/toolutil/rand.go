package toolutil

import (
	"math/rand"
	"time"
)

const strList string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStr(length int) string {
	var result []byte

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for i := 0; i < length; i++ {
		result = append(result, strList[r.Int63()%int64(len(strList))])
	}
	return string(result)
}

func RandInt(length int) int {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	var code int
	for i := 0; i < length; i++ {
		code += r.Intn(10)
	}
	return code
}
