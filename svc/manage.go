package svc

import (
	"database/sql"
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/constant"
	"github.com/ZhuoshanLuo/entry_task/dao"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/tool"
	"github.com/gin-gonic/gin"
	"regexp"
	"runtime/debug"
)

var (
	ManageRegister   = NewHandlerFunc(doManageRegister)
	ManageLogin      = NewHandlerFunc(doManageLogin)
	AddActivity      = NewHandlerFunc(doAddActivity)
	DelActivity      = NewHandlerFunc(doDelActivity)
	AddActivityType  = NewHandlerFunc(doAddActivityType)
	EditActivity     = NewHandlerFunc(doEditActivity)
	ShowAllUsers     = NewHandlerFunc(doShowAllUsers)
	EditActivityType = NewHandlerFunc(doEditActivityType)
	DelActivityType  = NewHandlerFunc(doDelActivityType)
	ShowActivityType = NewHandlerFunc(doShowActivityType)
)

//验证是否是管理员身份
func CheckIdentity(sessionId uint) (uint, bool, error) {
	userId, err := dao.QueryUserId(sessionId)
	if err != nil {
		return 0, false, err
	}
	//在user表中查找权限
	isAdmin, err := dao.QueryUserAdmin(userId)
	return userId, isAdmin, err
}

//运营人员注册接口
func doManageRegister(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.RegisterMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	name, passwd, email, avatar := req.Name, req.Passwd, req.Email, req.Avatar

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, passwd)
	if name == "" || passwd == "" || email == "" || !isMatch1 || !isMatch2 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//插入数据到user表中
	var user model.User
	user.Name, user.Avatar, user.Email = name, avatar, email
	user.Passwd = tool.AddSalt(passwd)
	user.IsAdmin = true
	user.CreatedAt = tool.GetTimeNowUnix()
	userId, err := dao.InsertUser(user)

	//插入数据库时发生错误
	if err != nil {
		tool.ErrorPrintln("sql insert into user table", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	//注册成功
	tool.InfoPrintln("insert into user table", userId)
	return codes.OK, nil
}

//运营人员登陆接口
func doManageLogin(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.LoginMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	name, reqPasswd := req.Name, req.Passwd

	//传入的用户名和密码不能为空
	if name == "" || reqPasswd == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	reqPasswd = tool.AddSalt(reqPasswd) //密码加盐
	userId, passwd, err := dao.QueryUserIsExist(name)
	//用户不存在
	if err == sql.ErrNoRows {
		tool.ErrorPrintln("sql user not exist", 0, debug.Stack())
		return codes.UserNotExist, nil
	}
	//密码错误
	if reqPasswd != passwd {
		return codes.PassWordError, nil
	}

	//用户是否已经登陆
	sessionId, err := dao.QueryUserIsLogin(userId)
	if err != nil {
		tool.ErrorPrintln("sql query user is login error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	if sessionId != 0 {
		return codes.OK, sessionId
	}

	//插入session数据库
	var session model.Session
	session.Id = tool.CreateSessionId(userId)
	session.UserId = userId
	err = dao.InsertSession(session)
	//数据库插入时出错
	if err != nil {
		tool.ErrorPrintln("sql insert into session table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	//登陆成功，返回sessoin
	tool.InfoPrintln("insert into session table", session.UserId)
	return codes.OK, session.Id
}

//添加活动
func doAddActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.AddActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, typeId, title, content := req.SessionId, req.TypeId, req.Title, req.Content
	location, start, end := req.Location, req.Start, req.End

	//参数不能为空
	if sessionId == 0 || typeId == 0 || title == "" || content == "" || location == "" || start == 0 || end == 0 {
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

	//插入activity表中
	var actObj model.Activity
	actObj.TypeId, actObj.Title, actObj.Content = typeId, title, content
	actObj.Location, actObj.StartTime, actObj.EndTime = location, start, end
	err = dao.InsertActvity(actObj)
	if err != nil {
		tool.ErrorPrintln("sql insert into activity table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("insert into activity table", userId)
	return codes.OK, nil
}

//删除活动
func doDelActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.DelActivityRequest
	c.BindJSON(&req)
	sessionId, actId := req.SessionId, req.ActId

	if sessionId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	userId, isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		tool.ErrorPrintln("user request a manage api", userId, debug.Stack())
		return codes.Forbidden, nil
	}

	//删除活动
	err = dao.DeleteActivity(actId)
	if err != nil {
		return codes.MysqlError, nil
	}
	return codes.OK, nil
}

//编辑活动
func doEditActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.EditActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, id, typeId, title, content := req.SessionId, req.Id, req.TypeId, req.Title, req.Content
	location, start, end := req.Location, req.Start, req.End

	if id == 0 || sessionId == 0 {
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

	err = dao.UpdateActivity(id, typeId, title, content, location, start, end)
	if err != nil {
		tool.ErrorPrintln("sql update activity table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("update actitivity table", userId)
	return codes.OK, nil
}

//添加活动类型
func doAddActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.AddActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error!", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, TypeName := req.SessionId, req.TypeName

	//缺少参数
	if sessionId == 0 || TypeName == "" {
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

	var obj model.ActivityType
	obj.Name = TypeName
	err = dao.InsertActivityType(obj)
	if err != nil {
		tool.ErrorPrintln("sql insert activity type table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("insert  actitivty type table", userId)
	return codes.OK, nil
}

//显示所有活动类型
func doShowActivityType(c *gin.Context) (codes.Code, interface{}) {
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

	rows, err := dao.QueryAllActivityType()
	if err != nil {
		tool.ErrorPrintln("sql query all activity type", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	var objects []model.ShowActivityTypeResponse
	for rows.Next() {
		var obj model.ShowActivityTypeResponse
		err := rows.Scan(&obj.TypeName)
		if err != nil {
			tool.ErrorPrintln("sql scan type name", userId, debug.Stack())
			return codes.MysqlError, nil
		}
		objects = append(objects, obj)
	}

	return codes.OK, objects
}

//编辑活动类型
func doEditActivityType(c *gin.Context) (codes.Code, interface{}) {
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

	err = dao.UpdateActivityType(typeId, name)
	if err != nil {
		tool.ErrorPrintln("sql update activity type", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("update activity type", userId)
	return codes.OK, nil
}

//删除活动类型
func doDelActivityType(c *gin.Context) (codes.Code, interface{}) {
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

	//删除活动类型
	err = dao.DeleteActivityType(typeId)
	if err != nil {
		tool.ErrorPrintln("sql delete activity type", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	tool.InfoPrintln("delete activity type", userId)
	return codes.OK, nil
}

//显示所有用户
func doShowAllUsers(c *gin.Context) (codes.Code, interface{}) {
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

	rows, err := dao.QueryAllUserMsg()
	if err != nil {
		tool.ErrorPrintln("sql query all user msg", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	var objects []model.UserPublicMsg
	for rows.Next() {
		var obj model.UserPublicMsg
		err := rows.Scan(&obj.Name, &obj.Email, &obj.Avatar)
		if err != nil {
			tool.ErrorPrintln("sql scan user msg error", userId, debug.Stack())
			return codes.MysqlError, nil
		}

		objects = append(objects, obj)
	}

	return codes.OK, objects
}
