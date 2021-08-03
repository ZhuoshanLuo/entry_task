package main

import (
	"database/sql"
	"example.com/greetings/globalVariable"
	"example.com/greetings/handler"
	"example.com/greetings/model"
	"fmt"
	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()

	var conf model.Config
	handler.GetConf(&conf)
	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", conf.Db.SqlUser, conf.Db.Passwd, conf.Db.Host, conf.Db.Database)
	globalVariable.DB, _ = sql.Open(conf.Db.Driver, sqlStr)
	globalVariable.DB.SetConnMaxLifetime(200)
	globalVariable.DB.SetMaxIdleConns(10)
	defer globalVariable.DB.Close()

	r.POST("/api/login", handler.Login)
	r.POST("/api/register", handler.Register)
	r.POST("/api/show_activities", handler.ShowActivities)
	r.POST("/api/activities_selector", handler.ActivitiesSelector)
	r.POST("/api/activity_info", handler.ActivityInfo)
	r.POST("/api/create_comment", handler.CreateComment)
	r.POST("/api/joined_activities_view", handler.ShowJoinedActivities)
	r.POST("/api/comments", handler.CommentsList)
	r.POST("/api/join_quit", handler.JoinOrExit)
	r.POST("api/activity_user_list", handler.ActivityUserList)
	r.POST("/manage/login", handler.Manage_login)
	r.POST("/manage/add_activity", handler.Manage_add_activity)
	r.POST("/manage_del_activity", handler.Manage_del_activity)
	r.POST("/manage/add_activity_type", handler.Manage_add_activity_type)
	r.POST("/manage/edit_activity", handler.Manage_edit_activity)
	r.POST("/manage/show_users", handler.Manage_show_users)
	r.POST("/manage/edit_activity_type", handler.Manage_edit_activity_type)
	r.POST("/manage/del_activity_type", handler.Manage_del_activity_type)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

