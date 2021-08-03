package main

import (
	"database/sql"
	"example.com/greetings/globalVariable"
	"example.com/greetings/model"
	"example.com/greetings/svc"
	"fmt"
	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()

	var conf model.Config
	svc.GetConf(&conf)
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.Db.SqlUser, conf.Db.Passwd, conf.Db.Host, conf.Db.Database)
	globalVariable.DB, _ = sql.Open(conf.Db.Driver, sqlStr)
	globalVariable.DB.SetConnMaxLifetime(200)
	globalVariable.DB.SetMaxIdleConns(10)
	defer globalVariable.DB.Close()

	r.POST("/api/login", svc.Login)
	r.POST("/api/register", svc.Register)
	r.POST("/api/show_activities", svc.ShowActivities)
	r.POST("/api/activities_selector", svc.ActivitiesSelector)
	r.POST("/api/activity_info", svc.ActivityInfo)
	r.POST("/api/create_comment", svc.CreateComment)
	r.POST("/api/joined_activities_view", svc.ShowJoinedActivities)
	r.POST("/api/join_quit", svc.JoinOrExit)
	r.POST("/manage/login", svc.ManageLogin)
	r.POST("/manage/add_activity", svc.AddActivity)
	r.POST("/manage_del_activity", svc.DelActivity)
	r.POST("/manage/add_activity_type", svc.AddActivityType)
	r.POST("/manage/edit_activity", svc.EditActivity)
	r.POST("/manage/show_users", svc.ShowAllUsers)
	r.POST("/manage/edit_activity_type", svc.EditActivityType)
	r.POST("/manage/del_activity_type", svc.DelActivityType)
	r.POST("/manage/del_activity_type", svc.ShowActivityType)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

