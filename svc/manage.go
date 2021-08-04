package svc

import (
	"database/sql"
	"example.com/greetings/codes"
	"example.com/greetings/constant"
	"example.com/greetings/dao"
	"example.com/greetings/model"
	"github.com/gin-gonic/gin"
	"regexp"
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

//运营人员注册接口
func doManageRegister(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.RegisterRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	name, passwd, email, avatar := req.Name, req.Passwd, req.Email, req.Avatar

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, passwd)
	if name == "" || passwd == "" || email == "" || !isMatch1 || !isMatch2 {
		return codes.MissParameter, nil
	}

	//用户已经存在
	isRegister, err := dao.QueryUserIsRegister(name, email)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isRegister {
		return codes.UserExist, nil
	}

	//插入数据到user表中
	var user model.User
	user.Name, user.Avatar, user.Email = name, avatar, email
	user.Passwd = AddSalt(passwd)
	user.IsAdmin = true
	user.CreatedAt = GetTime()
	err = dao.InsertUser(user)

	//插入数据库时发生错误
	if err != nil {
		return codes.MysqlError, nil
	}

	//注册成功
	return codes.OK, nil
}

//运营人员登陆接口
func doManageLogin(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.LoginRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	name, reqPasswd := req.Name, req.Passwd

	//传入的用户名和密码不能为空
	if name == "" || reqPasswd == "" {
		return codes.MissParameter, nil
	}

	reqPasswd = AddSalt(reqPasswd) //密码加盐
	userId, passwd, err := dao.QueryUserIsExist(name)
	//用户不存在
	if err == sql.ErrNoRows {
		return codes.UserNotExist, nil
	}
	//密码错误
	if reqPasswd != passwd {
		return codes.PassWordError, nil
	}

	//插入session数据库
	var session model.Session
	session.Id = CreateSessionId(userId)
	session.UserId = userId
	err = dao.InsertSession(session)
	//数据库插入时出错
	if err != nil {
		return codes.MysqlError, nil
	}
	//登陆成功，返回sessoin
	return codes.OK, session.Id
}

//添加活动
func doAddActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.AddActivityRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	sessionId, typeId, title, content, location, start, end := req.SessionId, req.TypeId, req.Title, req.Content, req.Location, req.Start, req.End

	//参数不能为空
	if sessionId == 0 || typeId == 0 || title == "" || content == "" || location == "" || start == 0 || end == 0 {
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	//插入activity表中
	var actObj model.Activity
	actObj.TypeId, actObj.Title, actObj.Content, actObj.Location, actObj.Start, actObj.End = typeId, title, content, location, start, end
	err = dao.InsertActvity(actObj)
	if err != nil {
		return codes.MysqlError, nil
	}

	return codes.OK, nil
}

//删除活动
func doDelActivity(c *gin.Context) (codes.Code, interface{}) {
	var req model.DelActivityRequest
	c.BindJSON(&req)
	sessionId, actId := req.SessionId, req.ActId

	if sessionId == 0 {
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
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
		return codes.BindJsonError, nil
	}
	sessionId, id, typeId, title, content, location, start, end := req.SessionId, req.Id, req.TypeId, req.Title, req.Content, req.Location, req.Start, req.End

	if id == 0 || sessionId == 0 {
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	err = dao.UpdateActivity(id, typeId, title, content, location, start, end)
	if err != nil {
		return codes.MysqlError, nil
	}

	return codes.OK, nil
}

//添加活动类型
func doAddActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.AddActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	sessionId, TypeName := req.SessionId, req.TypeName

	//缺少参数
	if sessionId == 0 || TypeName == "" {
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	var obj model.ActivityType
	obj.Name = TypeName
	err = dao.InsertActivityType(obj)
	if err != nil {
		return codes.MysqlError, nil
	}
	return codes.OK, nil
}

func doShowActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.ShowActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	sessionId := req.SessoinId

	if sessionId == 0 {
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	rows, err := dao.QueryAllActivityType()
	if err != nil {
		return codes.MysqlError, nil
	}

	var objects []model.ShowActivityTypeResponse
	for rows.Next() {
		var obj model.ShowActivityTypeResponse
		err := rows.Scan(&obj.TypeName)
		if err != nil {
			return codes.MysqlError, nil
		}
		objects = append(objects, obj)
	}
	return codes.OK, objects
}

func doEditActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.EditActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	sessionId, typeId, name := req.SessionId, req.Id, req.Name
	if sessionId == 0 || typeId == 0 || name == ""{
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	err = dao.UpdateActivityType(typeId, name)
	if err != nil {
		return codes.MysqlError, nil
	}

	return codes.OK, nil
}

func doDelActivityType(c *gin.Context) (codes.Code, interface{}) {
	var req model.DelActivityTypeRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	sessionId, typeId := req.SessionId, req.Id

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err == sql.ErrNoRows {
		return codes.Forbidden, nil
	}
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	//不能缺少参数
	if typeId == 0 {
		return codes.MissParameter, nil
	}

	//删除活动类型
	err = dao.DeleteActivityType(typeId)
	if err != nil {
		return codes.MysqlError, nil
	}
	return codes.OK, nil
}

func doShowAllUsers(c *gin.Context) (codes.Code, interface{}) {
	var req model.ShowAllUsersRequest
	err := c.BindJSON(&req)
	if err != nil {
		return codes.BindJsonError, nil
	}
	sessionId := req.SessionId

	if sessionId == 0 {
		return codes.MissParameter, nil
	}

	//首先先验证身份，是否已登陆，是否是运营人员
	isAdmin, err := CheckIdentity(sessionId)
	if err != nil {
		return codes.MysqlError, nil
	}
	if isAdmin == false {
		return codes.Forbidden, nil
	}

	rows, err := dao.QueryAllUserMsg()
	if err != nil {
		return codes.MysqlError, nil
	}

	var objects []model.ShowAllUsersResponse
	for rows.Next() {
		var obj model.ShowAllUsersResponse
		err := rows.Scan(&obj.Name, &obj.Email, &obj.Avatart)
		if err != nil {
			return codes.MysqlError, nil
		}

		objects = append(objects, obj)
	}
	return codes.OK, objects
}
