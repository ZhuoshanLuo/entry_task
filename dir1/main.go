package dir1

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
)

//常量
const (
	SqlUser     = "root"
	Passwd   = "Luo2566288"
	Host     = "127.0.0.1:3306"
	Database = "et_db"
	Salt     = "abcd"
)

//全局变量
var (
	db        *sql.DB
	BytesKind = reflect.TypeOf(sql.RawBytes{}).Kind()
	TimeKind  = reflect.TypeOf(mysql.NullTime{}).Kind()
)

//结构体
/*
type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Activity struct {
	Title string `json:"title"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
}

type ActivityDetail struct {
	Content  string   `json:"content"`
	Location string   `json:"location"`
	Act      Activity `json:"activity"`
}

type User struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	Email string `json:"email"`
}

type Comment struct {
	UserName string `json:"username"`
	Time     uint   `json:"time"`
	Content  string `json:"content"`
}

//活动选择器函数的返回值
type SelectorStruct struct {
	act Activity `json:"activity"`
	Res Response `json:"response"`
	JoinStatus uint `json:"status"'`
}

//活动详情的返回值
type ADetailStruct struct {

}
*/

//结构体

//type Activity struct {
//	Id            uint   `json:"id"`
//	Type_id       uint   `json:"type_id"`
//	Type_name     string `json:"type_name"`
//	Title         string `json:"title"`
//	Content       string `json:"content"`
//	Start_time    uint   `json:"start_time"`
//	End_time      uint   `json:"end_time"`
//	Activity_type string `json:"activity_type"`
//	Join_status   uint   `json:"join_status"`
//}

type UserMsg struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type Activity_detail struct {
	Title       string    `json:"title"`
	Start_time  uint      `json:"start_time"`
	End_time    uint      `json:"end_time"`
	Content     string    `json:"content"`
	Location    string    `json:"location"`
	Join_status uint      `json:"join_status"`
	Users       []UserMsg `json:"users"`
}

type Activities_joinin struct {
	Title      string `json:"title"`
	Start_time uint   `json:"start_time"`
	End_time   uint   `json:"end_time"`
}

type Comment struct {
	User_name    string `json:"user_name"`
	Created_time uint   `json:"create_time"`
	Content      string `json:"content"`
}

type Activity_select struct {
	Type_name   string `json:"type_name"`
	Title       string `json:"title"`
	Start_time  uint   `json:"start_time"`
	End_time    uint   `json:"end_time"`
	Join_status uint   `json:"join_status"`
}

