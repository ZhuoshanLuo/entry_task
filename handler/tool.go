package handler

import (
	"crypto/md5"
	"encoding/hex"
	"example.com/greetings/dir1"
	"hash/crc32"
	"time"
)

func AddSalt(passwd string) string {
	passwd = passwd + dir1.Salt
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(passwd))
	passwd = hex.EncodeToString(md5Ctx.Sum(nil))
	return passwd
}

//生成各种表的id字段
func CreateId(s string) uint {
	hash_code := uint(crc32.ChecksumIEEE([]byte(s)))
	if hash_code >= 0 {
		return hash_code
	}
	if hash_code < 0 {
		return -hash_code
	}
	return 0
}

func GetTime() uint {
	return uint(time.Now().Unix())
}

