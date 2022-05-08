package tool

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

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	return rand.Intn(max-min) + min
}
