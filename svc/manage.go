package svc

import (
	"database/sql"
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/dao"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/tool"
	"runtime/debug"
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
func doManageRegister(req model.RegisterMsg) (codes.Code, interface{}) {
	name, passwd, email, avatar := req.Name, req.Passwd, req.Email, req.Avatar

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
func doManageLogin(req model.LoginMsg) (codes.Code, interface{}) {
	name, reqPasswd := req.Name, req.Passwd

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
	if err != sql.ErrNoRows && err != nil {
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
func doAddActivity(req model.AddActivityRequest, userId uint) (codes.Code, interface{}) {
	typeId, title, content := req.TypeId, req.Title, req.Content
	location, start, end := req.Location, req.Start, req.End

	//插入activity表中
	var actObj model.Activity
	actObj.TypeId, actObj.Title, actObj.Content = typeId, title, content
	actObj.Location, actObj.StartTime, actObj.EndTime = location, start, end
	err := dao.InsertActvity(actObj)
	if err != nil {
		tool.ErrorPrintln("sql insert into activity table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("insert into activity table", userId)
	return codes.OK, nil
}

//删除活动
func doDelActivity(req model.DelActivityRequest, userId uint) (codes.Code, interface{}) {
	actId := req.ActId

	//删除活动
	err := dao.DeleteActivity(actId)
	if err != nil {
		return codes.MysqlError, nil
	}
	return codes.OK, nil
}

//编辑活动
func doEditActivity(req model.EditActivityRequest, userId uint) (codes.Code, interface{}) {
	id, typeId, title, content := req.Id, req.TypeId, req.Title, req.Content
	location, start, end := req.Location, req.Start, req.End

	err := dao.UpdateActivity(id, typeId, title, content, location, start, end)
	if err != nil {
		tool.ErrorPrintln("sql update activity table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("update actitivity table", userId)
	return codes.OK, nil
}

//添加活动类型
func doAddActivityType(req model.AddActivityTypeRequest, userId uint) (codes.Code, interface{}) {
	TypeName := req.TypeName

	var obj model.ActivityType
	obj.Name = TypeName
	err := dao.InsertActivityType(obj)
	if err != nil {
		tool.ErrorPrintln("sql insert activity type table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("insert  actitivty type table", userId)
	return codes.OK, nil
}

//显示所有活动类型
func doShowActivityType(req model.SessionId, userId uint) (codes.Code, interface{}) {
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
func doEditActivityType(req model.EditActivityTypeRequest, userId uint) (codes.Code, interface{}) {
	typeId, name := req.Id, req.Name

	err := dao.UpdateActivityType(typeId, name)
	if err != nil {
		tool.ErrorPrintln("sql update activity type", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("update activity type", userId)
	return codes.OK, nil
}

//删除活动类型
func doDelActivityType(req model.DelActivityTypeRequest, userId uint) (codes.Code, interface{}) {
	typeId := req.Id

	//删除活动类型
	err := dao.DeleteActivityType(typeId)
	if err != nil {
		tool.ErrorPrintln("sql delete activity type", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	tool.InfoPrintln("delete activity type", userId)
	return codes.OK, nil
}

//显示所有用户
func doShowAllUsers(req model.SessionId, userId uint) (codes.Code, interface{}) {

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
