package svc

import (
	"database/sql"
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/dao"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/utils"
	"runtime/debug"
)

/*
//运营人员注册接口
发送：
{
    "name" : "zgmm",
    "passwd" : "12345678",
    "email" : "abcdefge",
    "avatar" : ""
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//注册
func doManageRegister(req model.RegisterMsg) model.Response {
	name, passwd, email, avatar := req.Name, req.Passwd, req.Email, req.Avatar

	//插入数据到user表中
	var user model.User
	user.Name, user.Avatar, user.Email = name, avatar, email
	user.Passwd = utils.AddSalt(passwd)
	user.IsAdmin = true
	user.CreatedAt = utils.GetTimeNowUnix()
	userId, err := dao.InsertUser(user)

	//插入数据库时发生错误
	if err != nil {
		utils.ErrorPrintln("sql insert into user table", 0, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	//注册成功
	utils.InfoPrintln("insert into user table", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//运营人员登陆接口
发送：
{
    "name" : "zgmm",
    "passwd" : "12345678"
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": 1281784366
}
*/
//登陆
func doManageLogin(req model.LoginMsg) model.Response {
	name, reqPasswd := req.Name, req.Passwd

	reqPasswd = utils.AddSalt(reqPasswd) //密码加盐
	userId, passwd, err := dao.QueryUserIsExist(name)
	//用户不存在
	if err == sql.ErrNoRows {
		utils.ErrorPrintln("sql user not exist", 0, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.UserNotExist,
				Msg:  codes.Errorf(codes.UserNotExist),
			},
		}
	}
	//密码错误
	if reqPasswd != passwd {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.PassWordError,
				Msg:  codes.Errorf(codes.PassWordError),
			},
		}
	}

	//用户是否已经登陆
	sessionId, err := dao.QueryUserIsLogin(userId)
	if err != sql.ErrNoRows && err != nil {
		utils.ErrorPrintln("sql query user is login error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	if sessionId != 0 {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.OK,
				Msg:  codes.Errorf(codes.OK),
			},
			Session_id: sessionId,
		}
	}

	//插入session数据库
	var session model.Session
	session.Id = utils.CreateSessionId(userId)
	session.UserId = userId
	err = dao.InsertSession(session)
	//数据库插入时出错
	if err != nil {
		utils.ErrorPrintln("sql insert into session table error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	//登陆成功，返回sessoin
	utils.InfoPrintln("insert into session table", session.UserId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
		Session_id: sessionId,
	}
}

/*
//添加活动
发送：
{
    "type_id" : 1,
    "title" : "b",
    "content" : "test",
    "location" : "b",
    "start" : 2,
    "end" : 2
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//添加活动
func doAddActivity(req model.AddActivityRequest, userId uint) model.Response {
	typeId, title, content := req.TypeId, req.Title, req.Content
	location, start, end := req.Location, req.Start, req.End

	//插入activity表中
	var actObj model.Activity
	actObj.TypeId, actObj.Title, actObj.Content = typeId, title, content
	actObj.Location, actObj.StartTime, actObj.EndTime = location, start, end
	err := dao.InsertActvity(actObj)
	if err != nil {
		utils.ErrorPrintln("sql insert into activity table error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	utils.InfoPrintln("insert into activity table", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//删除活动
发送：
{
    "activityId" : 2
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//删除活动
func doDelActivity(req model.DelActivityRequest, userId uint) model.Response {
	actId := req.ActId

	//删除活动
	err := dao.DeleteActivity(actId)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//编辑活动
发送：
{
    "id" : 1,
    "content" : "test",
    "typeId" : 2
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//编辑活动
func doEditActivity(req model.EditActivityRequest, userId uint) model.Response {
	id, typeId, title, content := req.Id, req.TypeId, req.Title, req.Content
	location, start, end := req.Location, req.Start, req.End

	err := dao.UpdateActivity(id, typeId, title, content, location, start, end)
	if err != nil {
		utils.ErrorPrintln("sql update activity table error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	utils.InfoPrintln("update actitivity table", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//添加活动类型
发送：
{
    "type_name" : "sport"
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//添加活动类型
func doAddActivityType(req model.AddActivityTypeRequest, userId uint) model.Response {
	TypeName := req.TypeName

	var obj model.ActivityType
	obj.Name = TypeName
	err := dao.InsertActivityType(obj)
	if err != nil {
		utils.ErrorPrintln("sql insert activity type table error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	utils.InfoPrintln("insert  actitivty type table", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//显示所有活动类型
发送：
{
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": [
        {
            "typeName": "typeOne"
        },
        {
            "typeName": "typeTwo"
        }
    ]
}
*/
//显示所有活动类型
func doShowActivityType(req model.ShowActivityTypeRequest, userId uint) model.Response {
	page := req.Page
	rows, err := dao.QueryAllActivityType(page)
	if err != nil {
		utils.ErrorPrintln("sql query all activity type", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	var objects []model.ShowActivityTypeResponse
	for rows.Next() {
		var obj model.ShowActivityTypeResponse
		err := rows.Scan(&obj.TypeName)
		if err != nil {
			utils.ErrorPrintln("sql scan type name", userId, debug.Stack())
			return model.Response{
				Status: model.CodeMsg{
					Code: codes.MysqlError,
					Msg:  codes.Errorf(codes.MysqlError),
				},
			}
		}
		objects = append(objects, obj)
	}

	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
		Data: objects,
	}
}

/*
//编辑活动类型
发送：
{
    "id" : 4,
    "name" : "typeFour"
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//编辑活动类型
func doEditActivityType(req model.EditActivityTypeRequest, userId uint) model.Response {
	typeId, name := req.Id, req.Name

	err := dao.UpdateActivityType(typeId, name)
	if err != nil {
		utils.ErrorPrintln("sql update activity type", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	utils.InfoPrintln("update activity type", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//删除活动类型
发送：
{
    "id" : 4
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//删除活动类型
func doDelActivityType(req model.DelActivityTypeRequest, userId uint) model.Response {
	typeId := req.Id

	//删除活动类型
	err := dao.DeleteActivityType(typeId)
	if err != nil {
		utils.ErrorPrintln("sql delete activity type", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	utils.InfoPrintln("delete activity type", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//显示所有用户
发送：
{
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": [
        {
            "name": "lfp",
            "email": "123456",
            "avatar": ""
        },
        {
            "name": "lzs",
            "email": "lzs@shopee.com",
            "avatar": ""
        }
    ]
}
*/
//显示所有用户
func doShowAllUsers(req model.ShowAllUsersRequest, userId uint) model.Response {
	page := req.Page
	rows, err := dao.QueryAllUserMsg(page)
	if err != nil {
		utils.ErrorPrintln("sql query all user msg", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	var objects []model.UserPublicMsg
	for rows.Next() {
		var obj model.UserPublicMsg
		err := rows.Scan(&obj.Name, &obj.Email, &obj.Avatar)
		if err != nil {
			utils.ErrorPrintln("sql scan user msg error", userId, debug.Stack())
			return model.Response{
				Status: model.CodeMsg{
					Code: codes.MysqlError,
					Msg:  codes.Errorf(codes.MysqlError),
				},
			}
		}

		objects = append(objects, obj)
	}
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
		Data: objects,
	}
}
