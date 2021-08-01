package handler

import (
	"crypto/md5"
	"encoding/hex"
	"example.com/greetings/codes"
	"example.com/greetings/dir1"
	"example.com/greetings/model"
	"gopkg.in/yaml.v2"
	"hash/crc32"
	"io/ioutil"
	"log"
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
func CreateSessionId(s string) uint {
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

func ResponseFun(code codes.Code, data interface{}) model.Response {
	res := model.Response{
		Code: code,
		Msg:  codes.Errorf(code),
	}
	if data != nil {
		res.Data = data
	}
	return res
}

func GetConf(c *model.Config) {
	yamlFile, err := ioutil.ReadFile("config/server.yaml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

