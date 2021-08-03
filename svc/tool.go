package svc

import (
	"crypto/md5"
	"encoding/hex"
	"example.com/greetings/codes"
	"example.com/greetings/constant"
	"example.com/greetings/dao"
	"example.com/greetings/model"
	"fmt"
	"gopkg.in/yaml.v2"
	"hash/crc32"
	"io/ioutil"
	"log"
	"time"
)

func AddSalt(passwd string) string {
	passwd += constant.Salt
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(passwd))
	passwd = hex.EncodeToString(md5Ctx.Sum(nil))
	return passwd
}

//生成session表的id字段
func CreateSessionId(s uint) uint {
	id := string(s)
	hash_code := uint(crc32.ChecksumIEEE([]byte(id)))
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
	yamlFile, err := ioutil.ReadFile("./config/server.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	fmt.Println(c.Db.Driver, c.Db.Passwd)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
}

func CheckIdentity(sessionId uint) (bool, error) {
	userId, err := dao.QueryUserId(sessionId)
	if err != nil {
		return false, err
	}
	//在user表中查找权限
	return dao.QueryUserAdmin(userId)
}
