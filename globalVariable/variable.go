package globalVariable

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"reflect"
)

//全局变量
var (
	DB        *sql.DB
	BytesKind = reflect.TypeOf(sql.RawBytes{}).Kind()
	TimeKind  = reflect.TypeOf(mysql.NullTime{}).Kind()
)
