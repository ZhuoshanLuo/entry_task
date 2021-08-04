package globalVariable

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"reflect"
)

//全局变量
var (
	DB        *sqlx.DB
	BytesKind = reflect.TypeOf(sql.RawBytes{}).Kind()
	TimeKind  = reflect.TypeOf(mysql.NullTime{}).Kind()
)
