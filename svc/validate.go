package svc

import (
	"github.com/ZhuoshanLuo/entry_task/constant"
	"github.com/ZhuoshanLuo/entry_task/dao"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/utils"
	"github.com/gin-gonic/gin"
	"regexp"
	"runtime/debug"
	"strconv"
)

//注册接口参数处理
func registerParameterValid(c *gin.Context) (bool, *model.RegisterMsg) {
	//提取请求参数
	var req model.RegisterMsg
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, req.Name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, req.Email)
	if req.Name == "" || req.Passwd == "" || req.Email == "" || !isMatch1 || !isMatch2 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//登陆接口参数处理
func loginParameterValid(c *gin.Context) (bool, *model.LoginMsg) {
	//提取请求参数
	var req model.LoginMsg
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	//传入的用户名和密码不能为空
	if req.Name == "" || req.Passwd == "" {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//显示所有活动接口参数处理
func showActivitiesParameterValid(c *gin.Context) (bool, *model.ShowActivtyRequest) {
	var req model.ShowActivtyRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	//分页显示
	if req.Page == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//活动过滤接口参数处理
func activitiesSelectorParameterValid(c *gin.Context) (bool, *model.ActivitySelectorRequest) {
	var req model.ActivitySelectorRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}
	return true, &req
}

//发表评论参数处理
func createCommentParameterValid(c *gin.Context) (bool, *model.CreateCommentRequest) {
	var req model.CreateCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	//登陆后才能评论，评论的活动和内容不为空
	if req.ActivityId == 0 || req.Content == "" {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//用户加入或退出活动接口参数处理
func joinOrExitParameterValid(c *gin.Context) (bool, *model.JoinOrExitRequest) {
	var req model.JoinOrExitRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	//传入参数都不能为空或格式错误
	if req.ActivityId == 0 || req.Action == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//活动信息接口参数处理
func activityInfoParameterValid(c *gin.Context) (bool, *model.ActivityInfoRequest) {
	var req model.ActivityInfoRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	//活动id是必要的参数
	if req.ActivityId == 0 || req.CommentPage == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//运营人员注册接口参数处理
func manageRegisterParameterValid(c *gin.Context) (bool, *model.RegisterMsg) {
	//提取请求参数
	var req model.RegisterMsg
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, req.Name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, req.Passwd)
	if req.Name == "" || req.Passwd == "" || req.Email == "" || !isMatch1 || !isMatch2 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//运营人员登陆接口参数处理
func manageLoginParameterValid(c *gin.Context) (bool, *model.LoginMsg) {
	//提取请求参数
	var req model.LoginMsg
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}

	//传入的用户名和密码不能为空
	if req.Name == "" || req.Passwd == "" {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//添加活动接口参数处理
func addActivityParameterValid(c *gin.Context) (bool, *model.AddActivityRequest) {
	var req model.AddActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}

	//参数不能为空
	if req.TypeId == 0 || req.Title == "" || req.Content == "" || req.Location == "" || req.Start == 0 || req.End == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//删除活动接口参数处理
func delActivityParameterValid(c *gin.Context) (bool, *model.DelActivityRequest) {
	var req model.DelActivityRequest
	c.BindJSON(&req)

	return true, &req
}

//编辑活动接口参数处理
func editActivityParameterValid(c *gin.Context) (bool, *model.EditActivityRequest) {
	var req model.EditActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}

	//编辑的活动id不能为空
	if req.Id == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//添加活动类型接口参数处理
func addActivityTypeParameterValid(c *gin.Context) (bool, *model.AddActivityTypeRequest) {
	var req model.AddActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}

	//编辑的活动类型不能为空
	if req.TypeName == "" {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//编辑活动类型接口参数处理
func editActivityTypeParameterValid(c *gin.Context) (bool, *model.EditActivityTypeRequest) {
	var req model.EditActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}
	typeId, name := req.Id, req.Name

	if typeId == 0 || name == "" {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//删除活动类型接口参数处理
func delActivityTypeParameterValid(c *gin.Context) (bool, *model.DelActivityTypeRequest) {
	var req model.DelActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error!", 0, debug.Stack())
		return false, nil
	}
	typeId := req.Id

	//不能缺少参数
	if typeId == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//显示所有活动类型接口参数处理
func ShowActivityTypeRequestParameterValid(c *gin.Context) (bool, *model.ShowActivityTypeRequest) {
	var req model.ShowActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	if req.Page == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//显示所有用户接口的参数校验
func showAllUsersParameterValid(c *gin.Context) (bool, *model.ShowAllUsersRequest) {
	var req model.ShowAllUsersRequest
	err := c.BindJSON(&req)
	if err != nil {
		utils.ErrorPrintln("bind json error", 0, debug.Stack())
		return false, nil
	}

	if req.Page == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return false, nil
	}
	return true, &req
}

//验证是否是管理员身份
func CheckIdentity(c *gin.Context) (uint, bool) {
	sessionIdStr, err := c.Cookie("session_id")
	if err != nil {
		return 0, false
	}
	sessionIdint, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		return 0, false
	}
	sessionId := uint(sessionIdint)
	userId, err := dao.QueryUserId(sessionId)
	if err != nil {
		return 0, false
	}
	//在user表中查找权限
	isAdmin, err := dao.QueryUserAdmin(userId)
	if err != nil {
		return 0, false
	}
	if isAdmin == false {
		return 0, false
	}
	return userId, true
}

//检查用户是否登陆，返回user id和程序异常error
func CheckUserLogin(c *gin.Context) (uint, error) {
	sessionIdStr, err := c.Cookie("session_id")
	if err != nil {
		return 0, nil
	}
	sessionIdint, err := strconv.Atoi(sessionIdStr)
	if err != nil {
		return 0, err
	}
	sessionId := uint(sessionIdint)
	userId, err := dao.QueryUserId(sessionId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
