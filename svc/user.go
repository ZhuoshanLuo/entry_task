package svc

import (
	"database/sql"
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/dao"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/ZhuoshanLuo/entry_task/utils"
	"github.com/gin-gonic/gin"
	"runtime/debug"
)

type Handler func(c *gin.Context) model.Response

/*
//用户注册接口
发送：
{
    "name" : "lss",
    "passwd" : "hahahaha",
    "email" : "123456789",
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
func doRegister(req model.RegisterMsg) model.Response {
	name, passwd, email, avatar := req.Name, req.Passwd, req.Email, req.Avatar

	//插入数据到user表中
	var user model.User
	user.Name, user.Avatar, user.Email = name, avatar, email
	user.Passwd = utils.AddSalt(passwd)
	user.IsAdmin = false
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
//用户登陆接口,允许重复登陆
发送：
{
    "name" : "lss",
    "passwd" : "hahahaha"
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": 996231864
}
*/
//登陆
func doLogin(req model.LoginMsg) model.Response {
	name, passwd := req.Name, req.Passwd

	passwd = utils.AddSalt(passwd) //密码加盐
	userId, passwd, err := dao.QueryUserIsExist(name)
	//用户不存在
	if err == sql.ErrNoRows {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.UserNotExist,
				Msg:  codes.Errorf(codes.UserExist),
			},
		}
	}
	//操作数据库出错
	if err != nil {
		utils.ErrorPrintln("sql query user is login error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	//密码错误
	if passwd != passwd {
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
		utils.InfoPrintln("user success login", sessionId)
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
	utils.InfoPrintln("user success login", session.UserId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
		Session_id: sessionId,
	}
}

/*
//显示所有活动
发送：
{
    "page" : 0
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": {
        "activities": [
            {
                "activity_profile": {
                    "title": "a",
                    "start": 1,
                    "end": 1
                },
                "join_status": true
            }
        ]
    }
}
*/
//显示所有活动
func doShowActivities(req model.ShowActivtyRequest, userId uint) model.Response {
	page := req.Page

	//如果用户是登陆状态，传入sessionId不为空
	var err error

	//获取活动表的当前page
	rows, err := dao.QueryALLActivityRows(page)
	//查询时出错
	if err != nil {
		utils.ErrorPrintln("sql query activity table error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	//构建返回数据
	var objects model.ShowActivitiesResponse
	for rows.Next() {
		var actId uint
		var obj model.UserActivityInfo
		rows.Scan(&actId, &obj.Title, &obj.Start, &obj.End)

		if userId != 0 {
			obj.JoinStatus, err = dao.IsJoinin(userId, actId)
			if err != nil {
				utils.ErrorPrintln("sql query user is joinin from form table error", userId, debug.Stack())
				return model.Response{
					Status: model.CodeMsg{
						Code: codes.MysqlError,
						Msg:  codes.Errorf(codes.MysqlError),
					},
				}
			}
		} else {
			obj.JoinStatus = false
		}

		objects.Activities = append(objects.Activities, &obj)
	}

	//处理完成，返回活动列表objects
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
		Data: objects,
	}
}

/*
//活动过滤器
发送：
{
	//没有必要的参数
	//可选参数session_id保存在cookie中
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": {
        "activities": [
            {
                "activity_profile": {
                    "title": "a",
                    "start": 1,
                    "end": 1
                },
                "join_status": true
            }
        ]
    }
}
*/
//活动选择器
func doActivitiesSelector(req model.ActivitySelectorRequest, userId uint) model.Response {
	actType, start, end, page := req.Type, req.Start, req.End, req.Page

	//如果用户是登陆状态，获取user_id
	var err error

	//封装sql语句，得到符合条件的活动的行
	rows, err := dao.SqlActivitiesSelector(actType, start, end, page)
	if err != nil {
		utils.ErrorPrintln("sql query activities by condition error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	//遍历等到activity_selector的id、type_id、title、start_time、end_time
	var actId uint
	var objects model.ActivitySelectorResponse
	for rows.Next() {
		var obj model.UserActivityInfo
		rows.Scan(&actId, &obj.Title, &obj.Start, &obj.End)

		//是否参加活动
		if userId != 0 {
			obj.JoinStatus, err = dao.IsJoinin(userId, actId)
			if err != nil {
				utils.ErrorPrintln("sql query user is joinin from form table error", userId, debug.Stack())
				return model.Response{
					Status: model.CodeMsg{
						Code: codes.MysqlError,
						Msg:  codes.Errorf(codes.MysqlError),
					},
				}
			}
		} else {
			obj.JoinStatus = false
		}

		objects.Activities = append(objects.Activities, &obj)
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
//发表评论
发送：
{
    "activity_id" : 1,
    "content" : "good"
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//创建评论
func doCreateComment(req model.CreateCommentRequest, userId uint) model.Response {
	activityId, content := req.ActivityId, req.Content

	//将数据插入comment表中
	var comment model.Comment
	comment.UserId, comment.ActivityId, comment.Content = userId, activityId, content
	comment.CreatedAt = utils.GetTimeNowUnix()
	err := dao.InsertComment(comment)
	//插入数据库表时出错
	if err != nil {
		utils.ErrorPrintln("sql insert into comment table error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	utils.InfoPrintln("insert into comment table", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

/*
//显示用户加入的所有活动
发送：
{
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": {
        "activities": [
            {
                "title": "a",
                "start": 1,
                "end": 1
            }
        ]
    }
}
*/
//显示加入的所用活动
func doShowJoinedActivities(userId uint) model.Response {
	//获取用户加入的所有活动的id
	rows, err := dao.GetAllJoinActivities(userId)
	if err != nil {
		utils.ErrorPrintln("sql query all joinin activities error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	//获取活动的信息，标题、开始时间、结束时间
	var actId uint
	var objects model.ShowJoinedActivityResponse
	for rows.Next() {
		var obj model.ActivityProfile
		rows.Scan(&actId)

		obj.Title, obj.Start, obj.End, err = dao.QueryActivityMsg(actId)
		if err != nil {
			utils.ErrorPrintln("sql query activity msg error", userId, debug.Stack())
			return model.Response{
				Status: model.CodeMsg{
					Code: codes.MysqlError,
					Msg:  codes.Errorf(codes.MysqlError),
				},
			}
		}

		objects.Activities = append(objects.Activities, &obj)
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
//用户加入或退出活动
发送：
{
    "activity_id" : 1,
    "action" : 2  //1表示加入，2表示退出
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": null
}
*/
//加入或退出活动
func doJoinOrExit(req model.JoinOrExitRequest, userId uint) model.Response {
	actId, action := req.ActivityId, req.Action

	//在form表中插入或删除表项
	var err error
	if action == 2 {
		err = dao.DeleteForm(userId, actId)
	} else {
		var form model.Form
		form.ActId, form.UserId = actId, userId
		form.JoinedAt = utils.GetTimeNowUnix()
		err = dao.InsertForm(form)
	}
	//操作数据库时出错
	if err != nil {
		utils.ErrorPrintln("sql when delete or insert into form table", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	utils.InfoPrintln("user join or exit activity success", userId)
	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
	}
}

//显示活动详情的一部分：显示参加某个活动的所有用户
func GetActivityUserList(actId uint) (codes.Code, []model.UserPublicMsg) {
	//活动id是必须的
	if actId == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//从form表中获取所有加入活动的用户id
	rows, err := dao.QueryAllUsersId(actId)
	if err != nil {
		utils.ErrorPrintln("sql query user id error", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	//获取每个用户的信息
	var userId uint
	var objects []model.UserPublicMsg
	for rows.Next() {
		var obj model.UserPublicMsg
		err = rows.Scan(&userId)
		if err != nil {
			utils.ErrorPrintln("sql get user id error", 0, debug.Stack())
			return codes.MysqlError, nil
		}

		row := dao.QueryUsersMsg(userId)
		err := row.StructScan(&obj)
		if err != nil {
			utils.ErrorPrintln("sql query user msg error", 0, debug.Stack())
			return codes.MysqlError, nil
		}

		objects = append(objects, obj)
	}

	return codes.OK, objects
}

func GetCommentsList(actId uint, page uint) (codes.Code, []model.CommentListResponse) {
	//缺少活动的id参数
	if actId == 0 {
		utils.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//在comment_tab表中获取user_id、content、created_time
	rows, err := dao.QueryCommentMsg(actId, page)
	if err != nil {
		utils.ErrorPrintln("sql query comment msg", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	var userId uint
	var objects []model.CommentListResponse
	for rows.Next() {
		var obj model.CommentListResponse
		rows.Scan(&userId, &obj.Content, &obj.CreatedAt)

		obj.Name, err = dao.QueryUserName(userId)
		if err != nil {
			utils.ErrorPrintln("sql query user name error", 0, debug.Stack())
			return codes.MysqlError, nil
		}

		objects = append(objects, obj)
	}

	return codes.OK, objects
}

//获取活动详情
func GetActivityDetail(actId uint, userId uint) (*model.ActivityDetail, codes.Code) {
	//查询活动信息
	var actDetail model.ActivityDetail
	row := dao.QueryActivityDetail(actId)
	err := row.StructScan(&actDetail)
	if err != nil {
		utils.ErrorPrintln("sql query activity deltail error", userId, debug.Stack())
		return nil, codes.MysqlError
	}

	//用户是否加入此活动
	if userId != 0 {
		actDetail.JoinStatus, err = dao.IsJoinin(userId, actId)
		if err != nil {
			utils.ErrorPrintln("sql query user is joinin", userId, debug.Stack())
			return nil, codes.MysqlError
		}
	} else {
		actDetail.JoinStatus = false
	}
	return &actDetail, codes.OK
}

/*
//显示活动详情
发送：
{
    "activity_id" : 1,
    "comment_page" : 0
}
接收：
{
    "code": 0,
    "msg": "Success!",
    "data": {
        "Title": "a",
        "Start": 1,
        "End": 1,
        "Location": "a",
        "Content": "test",
        "joinStatus": true,
        "UserList": [
            {
                "Name": "lfp",
                "Email": "",
                "Avatar": ""
            }
        ],
        "CommentList": [
            {
                "name": "lfp",
                "content": "skr",
                "createdAt": 1627888219
            }
        ]
    }
}
*/
func doActivityInfo(req model.ActivityInfoRequest, userId uint) model.Response {
	actId, page := req.ActivityId, req.CommentPage

	var code codes.Code
	var obj model.ActivityInfoResponse
	obj.ActivityDetail, code = GetActivityDetail(actId, userId)
	if code != codes.OK {
		return model.Response{
			Status: model.CodeMsg{
				Code: code,
				Msg:  codes.Errorf(code),
			},
		}
	}

	//查询加入activity id的所有用户
	code, obj.UserList = GetActivityUserList(actId)
	if code != codes.OK {
		utils.ErrorPrintln("sql query user list error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: code,
				Msg:  codes.Errorf(code),
			},
		}
	}

	//查询某个活动的评论列表
	code, obj.CommentList = GetCommentsList(actId, page)
	if code != codes.OK {
		utils.ErrorPrintln("sql query comment list error", userId, debug.Stack())
		return model.Response{
			Status: model.CodeMsg{
				Code: code,
				Msg:  codes.Errorf(code),
			},
		}
	}

	return model.Response{
		Status: model.CodeMsg{
			Code: codes.OK,
			Msg:  codes.Errorf(codes.OK),
		},
		Data: obj,
	}
}
