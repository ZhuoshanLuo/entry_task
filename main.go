package main

import (
	"github.com/ZhuoshanLuo/entry_task/dao"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/svc"
	"github.com/ZhuoshanLuo/entry_task/tool"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	//读取配置参数
	var conf model.Config
	tool.GetConf(&conf)
	//连接数据库
	dao.Init(conf.Database)
	//打开日志
	tool.InitLog()

	//用户接口
	r.POST("/api/login", svc.Login)
	r.POST("/api/register", svc.Register)
	r.POST("/api/create_comment", svc.CreateComment)
	r.POST("/api/show_activities", svc.ShowActivities)
	r.POST("/api/activities_selector", svc.ActivitiesSelector)
	r.POST("/api/activity_info", svc.ActivityInfo)
	r.POST("/api/joined_activities_view", svc.ShowJoinedActivities)
	r.POST("/api/join_quit", svc.JoinOrExit)

	//运营后台接口
	r.POST("/manage/register", svc.ManageRegister)
	r.POST("/manage/login", svc.ManageLogin)
	//活动的增删查改
	r.POST("/manage/add_activity", svc.AddActivity)
	r.POST("/manage/del_activity", svc.DelActivity)
	r.POST("/manage/edit_activity", svc.EditActivity)
	r.POST("/manage/show_activities", svc.ShowActivities)
	//对活动类型的增删查改
	r.POST("/manage/show_activities_type", svc.ShowActivityType)
	r.POST("/manage/add_activity_type", svc.AddActivityType)
	r.POST("/manage/edit_activity_type", svc.EditActivityType)
	r.POST("/manage/del_activity_type", svc.DelActivityType)
	r.POST("/manage/show_users", svc.ShowAllUsers)

	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
