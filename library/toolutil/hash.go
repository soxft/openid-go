package toolutil

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha1(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Sha256(str, salt string) string {
	h := sha1.New()
	h.Write([]byte(str))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}
