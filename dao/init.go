package dao

import (
	"fmt"
	"github.com/ZhuoshanLuo/entry_task/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Init(db model.DatabaseConf) {
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", db.SqlUser, db.Passwd, db.Host, db.Database)
	DB, _ = sqlx.Open(db.Driver, sqlStr)
	DB.SetMaxOpenConns(200)
	DB.SetMaxIdleConns(10)
}
