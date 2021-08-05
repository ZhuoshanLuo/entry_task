package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/ZhuoshanLuo/entry_task/constant"
	"github.com/ZhuoshanLuo/entry_task/model"
	"gopkg.in/yaml.v2"
	"hash/crc32"
	"io/ioutil"
	"log"
	"runtime/debug"
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

func GetTimeNowUnix() uint {
	return uint(time.Now().Unix())
}

func GetConf(c *model.Config) {
	yamlFile, err := ioutil.ReadFile("./config/server.yml")
	if err != nil {
		FatalPrintln("Get yamlFile error", 0, debug.Stack())
	}

	err = yaml.Unmarshal(yamlFile, c)
	fmt.Println(c.Database.Driver, c.Database.Passwd)
	if err != nil {
		FatalPrintln("Unmarshal config", 0, debug.Stack())
		log.Fatalf("Unmarshal: %v", err)
	}
}
