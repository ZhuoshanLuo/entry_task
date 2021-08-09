package utils

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"math/rand"
	"strings"
	"time"
)

const (
	UserName      = "root"
	Password      = "Luo2566288"
	Ip            = "127.0.0.1"
	Port          = "3306"
	DBName        = "et_db"
	UserTable     = "user_tab"
	ActivityTable = "activities_tab"
)

const (
	USER_TOTAL_INSERT_NUM     = 1000000
	ACTIVITY_TOTAL_INSERT_NUM = 5000
	PER_INSERT_NUM            = 5000
	MAX_FAILNUM               = 0
)

const STRCHAR = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&*+-./:;<=>?@[]^_{|}~"

//随机生成一个整形数据
func MakeRandInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}

//随机生成一个字符串
func MakeRandString(minLength int64, maxLength int64) string {
	size := MakeRandInt(minLength, maxLength)
	str := make([]byte, size)
	for i := 0; i < int(size); i++ {
		index := MakeRandInt(0, int64(len(STRCHAR)))
		str[i] = STRCHAR[index]
	}
	return string(str)
}

func GetTime() int64 {
	return time.Now().Unix()
}

//连接数据库
func InitDB() *sql.DB {
	source := strings.Join([]string{UserName, ":", Password, "@tcp(", Ip, ":", Port, ")/", DBName, "?charset=utf8"}, "")
	DB, err := sql.Open("mysql", source)
	if err != nil {
		panic(fmt.Sprintf("Open Kungate Connection:[%s] failed, error is [%v].", source, err))
	}
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		fmt.Println("opon database fail")
		return nil
	}

	fmt.Println("connnect success")
	return DB
}

func UserInsert(db *sql.DB) {
	fmt.Printf("begin insert, total num:[%d]\n", USER_TOTAL_INSERT_NUM)
	var name string
	var passwd string
	var email string
	var avatar string
	var createdAt int64

	failnum := 0

	//拼接insert语句
	var strInsert = "insert into " + UserTable + "(name, passwd, email, avatar, is_admin, created_at) values"
	var InsertBuf = strInsert

	for i := 1; i <= USER_TOTAL_INSERT_NUM; i++ {
		name = MakeRandString(5, 50)
		passwd = MakeRandString(8, 20)
		email = MakeRandString(5, 50)
		avatar = MakeRandString(5, 50)
		createdAt = GetTime()

		InsertBuf += fmt.Sprintf(" ('%s', '%s', '%s', '%s', %t, %d)", name, passwd, email, avatar, false, createdAt)

		if i%PER_INSERT_NUM == 0 {
			InsertBuf += ";"

			_, err := db.Exec(InsertBuf)
			if err != nil {
				failnum++
				if failnum > MAX_FAILNUM {
					fmt.Println(err)
					return
				}

				InsertBuf = strInsert
			}

			InsertBuf = strInsert
		} else {
			InsertBuf += ","
		}
	}

}

func ActivityInsert(db *sql.DB) {
	fmt.Printf("begin insert, total num:[%d]\n", ACTIVITY_TOTAL_INSERT_NUM)
	var typeId int64
	var title string
	var content string
	var location string

	failnum := 0

	//拼接insert语句
	var strInsert = "insert into " + ActivityTable + "(type_id, title, content, location, start_time, end_time) values"
	var InsertBuf = strInsert

	for i := 1; i <= ACTIVITY_TOTAL_INSERT_NUM; i++ {
		typeId = MakeRandInt(1, 20)
		title = MakeRandString(5, 50)
		content = MakeRandString(5, 255)
		location = MakeRandString(5, 50)

		InsertBuf += fmt.Sprintf(" ('%d', '%s', '%s', '%s', %d, %d)", typeId, title, content, location, 1, 2)

		if i%PER_INSERT_NUM == 0 {
			InsertBuf += ";"

			_, err := db.Exec(InsertBuf)
			if err != nil {
				failnum++
				if failnum > MAX_FAILNUM {
					return
				}

				InsertBuf = strInsert
			}

			InsertBuf = strInsert
		} else {
			InsertBuf += ","
		}
	}
}

/*
func main() {
	db := InitDB()
	//UserInsert(db)
	ActivityInsert(db)
	db.Close()
}
*/
