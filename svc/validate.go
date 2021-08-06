package svc

import (
	"database/sql"
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/constant"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/tool"
	"github.com/gin-gonic/gin"
	"regexp"
	"runtime/debug"
)

/*
参数处理和权限处理文件
*/

//注册接口参数处理
func register(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.RegisterMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, req.Name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, req.Email)
	if req.Name == "" || req.Passwd == "" || req.Email == "" || !isMatch1 || !isMatch2 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	return doRegister(req)
}

//登陆接口参数处理
func login(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.LoginMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//传入的用户名和密码不能为空
	if req.Name == "" || req.Passwd == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doLogin(req)
}

//显示所有活动接口参数处理
func showActivities(c *gin.Context) (codes.Code, interface{}) {
	var req model.ShowActivtiyRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//分页显示
	if req.Page == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	return doShowActivities(req)
}

//活动过滤接口参数处理
func activitiesSelector(c *gin.Context) (codes.Code, interface{}) {
	var req model.ActivitySelectorRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	return doActivitiesSelector(req)
}

//发表评论参数处理
func createComment(c *gin.Context) (codes.Code, interface{}) {
	var req model.CreateCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//登陆后才能评论，评论的活动和内容不为空
	if req.SessionId == 0 || req.ActivityId == 0 || req.Content == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doCreateComment(req)
}

//显示用户加入的所有活动接口参数处理
func showJoinedActivities(c *gin.Context) (codes.Code, interface{}) {
	var req model.SessionId
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//用户需要在登陆状态，session_id不能为空
	if req.Id == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doShowJoinedActivities(req)
}

//用户加入或退出活动接口参数处理、权限管理（登陆）
func joinOrExit(c *gin.Context) (codes.Code, interface{}) {
	var req model.JoinOrExitRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//传入参数都不能为空或格式错误
	if req.SessionId == 0 || req.ActivityId == 0 || req.Action == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doJoinOrExit(req)
}

//活动信息接口参数处理
func activityInfo(c *gin.Context) (codes.Code, interface{}) {
	var req model.ActivityInfoRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//活动id是必要的参数
	if req.ActivityId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doActivityInfo(req)
}

//运营人员注册接口参数处理
func manageRegister(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.RegisterMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, req.Name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, req.Passwd)
	if req.Name == "" || req.Passwd == "" || req.Email == "" || !isMatch1 || !isMatch2 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doManageRegister(req)
}

//运营人员登陆接口接口参数处理
func manageLogin(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.LoginMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//传入的用户名和密码不能为空
	if req.Name == "" || req.Passwd == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}
	return doManageLogin(req)
}

//添加活动接口参数处理和管理员身份验证
func addActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.AddActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//参数不能为空
	if req.SessionId == 0 || req.TypeId == 0 || req.Title == "" || req.Content == "" || req.Location == "" || req.Start == 0 || req.End == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(req.SessionId)
	if err != nil {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doAddActivity(req, userId)
}

//删除活动接口参数处理和权限管理
func delActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.DelActivityRequest
	c.BindJSON(&req)

	if req.SessionId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(req.SessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doDelActivity(req, userId)
}

//编辑活动接口参数处理
func editActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.EditActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//编辑的活动id不能为空，管理人员需要在登陆状态
	if req.Id == 0 || req.SessionId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(req.SessionId)
	if err != nil {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doEditActivity(req, userId)
}

//添加活动类型接口参数处理和管理员身份验证
func addActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.AddActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}

	//session不能为空，编辑的活动类型不能为空
	if req.SessionId == 0 || req.TypeName == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(req.SessionId)
	if err != nil {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doAddActivityType(req, userId)
}

//显示活动类型接口参数处理和管理员身份验证
func showActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.SessionId
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId := req.Id

	if sessionId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doShowActivityType(req, userId)
}

//编辑活动类型接口参数处理和管理员身份验证
func editActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.EditActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, typeId, name := req.SessionId, req.Id, req.Name

	if sessionId == 0 || typeId == 0 || name == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doEditActivityType(req, userId)
}

//删除活动类型接口参数处理和管理员身份验证
func delActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.DelActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, typeId := req.SessionId, req.Id

	//不能缺少参数
	if typeId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(sessionId)
	if err == sql.ErrNoRows {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	if err != nil {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}
	return doDelActivityType(req, userId)
}

//显示所有活动接口参数处理和管理员身份验证
func showAllUsers(c *gin.Context) (codes.Code, interface{}) {
	var req model.SessionId
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId := req.Id

	if sessionId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		tool.ErrorPrintln("sql check user identity error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}
	return doShowAllUsers(req, userId)
}
