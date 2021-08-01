package main

import (
	"database/sql"
	"example.com/greetings/constant"
	"example.com/greetings/globalVariable"
	"example.com/greetings/handler"
	"fmt"
	"github.com/gin-gonic/gin"
)


func main() {
	r := gin.Default()

	sqlStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", constant.SqlUser, constant.Passwd, constant.Host, constant.Database)
	globalVariable.DB, _ = sql.Open(constant.Driver, sqlStr)
	globalVariable.DB.SetConnMaxLifetime(200)
	globalVariable.DB.SetMaxIdleConns(10)
	defer globalVariable.DB.Close()

	r.POST("/api/login", handler.Login)
	r.POST("/api/register", handler.Register)
	r.POST("/api/show_activities", handler.Show_activities)
	r.POST("/api/activities_selector", handler.Activities_selector)
	r.POST("/api/activity_info", handler.Activity_info)
	r.POST("/api/create_comment", handler.Create_comment)
	r.POST("/api/joined_activities_view", handler.Joined_activities_view)
	r.POST("/api/comments", handler.Comments)
	r.POST("/api/join_quit", handler.Join_quit)
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

