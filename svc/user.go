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

type Handler func(c *gin.Context) (codes.Code, interface{})

var (
	Register             = NewHandlerFunc(doRegister)
	Login                = NewHandlerFunc(doLogin)
	ShowActivities       = NewHandlerFunc(doShowActivities)
	ActivitiesSelector   = NewHandlerFunc(doActivitiesSelector)
	ActivityInfo         = NewHandlerFunc(doActivityInfo)
	CreateComment        = NewHandlerFunc(doCreateComment)
	ShowJoinedActivities = NewHandlerFunc(doShowJoinedActivities)
	JoinOrExit           = NewHandlerFunc(doJoinOrExit)
)

//用户注册接口
func doRegister(c *gin.Context) (codes.Code, interface{}) {
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
	user.IsAdmin = false
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

//用户登陆接口,允许重复登陆
func doLogin(c *gin.Context) (codes.Code, interface{}) {
	//提取请求参数
	var req model.LoginMsg
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	name, passwd := req.Name, req.Passwd

	//传入的用户名和密码不能为空
	if name == "" || passwd == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	passwd = tool.AddSalt(passwd) //密码加盐
	userId, passwd, err := dao.QueryUserIsExist(name)
	//用户不存在
	if err == sql.ErrNoRows {
		tool.ErrorPrintln("sql user not exist", 0, debug.Stack())
		return codes.UserNotExist, nil
	}
	//密码错误
	if passwd != passwd {
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

//显示所有活动
func doShowActivities(c *gin.Context) (codes.Code, interface{}) {
	var req model.ShowActivtiyRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, page := req.SessionId, req.Page

	//如果用户是登陆状态，传入sessionId不为空
	var userId uint
	if sessionId != 0 {
		userId, err = dao.QueryUserId(sessionId)
		//数据库执行时错误
		if err != nil {
			tool.ErrorPrintln("sql query user id from session table error", 0, debug.Stack())
			return codes.MysqlError, nil
		}
	}

	//获取活动表的当前page
	rows, err := dao.QueryALLActivityRows(page)
	//查询时出错
	if err != nil {
		tool.ErrorPrintln("sql query activity table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	//构建返回数据
	var objects model.ShowActivitiesResponse
	for rows.Next() {
		var actId uint
		var obj model.UserActivityInfo
		rows.Scan(&actId, &obj.Title, &obj.Start, &obj.End)

		if req.SessionId != 0 {
			obj.JoinStatus, err = dao.IsJoinin(userId, actId)
			if err != nil {
				tool.ErrorPrintln("sql query user is joinin from form table error", userId, debug.Stack())
				return codes.MysqlError, nil
			}
		} else {
			obj.JoinStatus = false
		}

		objects.Activities = append(objects.Activities, &obj)
	}

	//处理完成，返回活动列表objects
	return codes.OK, objects
}

//活动过滤器
func doActivitiesSelector(c *gin.Context) (codes.Code, interface{}) {
	var req model.ActivitySelectorRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, actType, start, end, page := req.SessionId, req.Type, req.Start, req.End, req.Page

	//如果用户是登陆状态，获取user_id
	var userId uint
	if sessionId != 0 {
		userId, err = dao.QueryUserId(sessionId)
		if err != nil {
			tool.ErrorPrintln("sql query user id error", 0, debug.Stack())
			return codes.MysqlError, nil
		}
	}

	//封装sql语句，得到符合条件的活动的行
	rows, err := dao.SqlActivitiesSelector(actType, start, end, page)
	if err != nil {
		tool.ErrorPrintln("sql query activities by condition error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	//遍历等到activity_selector的id、type_id、title、start_time、end_time
	var actId uint
	var objects model.ActivitySelectorResponse
	for rows.Next() {
		var obj model.UserActivityInfo
		rows.Scan(&actId, &obj.Title, &obj.Start, &obj.End)

		//是否参加活动
		if sessionId != 0 {
			obj.JoinStatus, err = dao.IsJoinin(userId, actId)
			if err != nil {
				tool.ErrorPrintln("sql query user is joinin from form table error", userId, debug.Stack())
				return codes.MysqlError, nil
			}
		} else {
			obj.JoinStatus = false
		}

		objects.Activities = append(objects.Activities, &obj)
	}
	return codes.OK, objects
}

//发表评论
func doCreateComment(c *gin.Context) (codes.Code, interface{}) {
	var req model.CreateCommentRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, activityId, content := req.SessionId, req.ActivityId, req.Content

	//传入内容不能为空
	if sessionId == 0 || activityId == 0 || content == "" {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//从session表中获取用户user id，查询不到session id对应的用户时也会报错
	userId, err := dao.QueryUserId(sessionId)
	//数据库查表时出错
	if err != nil {
		tool.ErrorPrintln("sql query user id from session table error", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	//将数据插入comment表中
	var comment model.Comment
	comment.UserId, comment.ActivityId, comment.Content = userId, activityId, content
	comment.CreatedAt = tool.GetTimeNowUnix()
	err = dao.InsertComment(comment)
	//插入数据库表时出错
	if err != nil {
		tool.ErrorPrintln("sql insert into comment table error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("insert into comment table", userId)
	return codes.OK, nil
}

//显示用户加入的所有活动
func doShowJoinedActivities(c *gin.Context) (codes.Code, interface{}) {
	var req model.SessionId
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId := req.Id

	//用户需要在登陆状态，session_id不能为空
	if sessionId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//获取用户user_id
	userId, err := dao.QueryUserId(sessionId)
	if err != nil {
		tool.ErrorPrintln("sql query user id from session table error", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	//获取用户加入的所有活动的id
	rows, err := dao.GetAllJoinActivities(userId)
	if err != nil {
		tool.ErrorPrintln("sql query all joinin activities error", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	//获取活动的信息，标题、开始时间、结束时间
	var actId uint
	var objects model.ShowJoinedActivityResponse
	for rows.Next() {
		var obj model.ActivityProfile
		rows.Scan(&actId)

		obj.Title, obj.Start, obj.End, err = dao.QueryActivityMsg(actId)
		if err != nil {
			tool.ErrorPrintln("sql query activity msg error", userId, debug.Stack())
			return codes.MysqlError, nil
		}

		objects.Activities = append(objects.Activities, &obj)
	}

	return codes.OK, objects
}

//用户加入或退出活动
func doJoinOrExit(c *gin.Context) (codes.Code, interface{}) {
	var req model.JoinOrExitRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, actId, action := req.SessionId, req.ActivityId, req.Action

	//传入参数都不能为空或格式错误
	if sessionId == 0 || actId == 0 || action == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//查看用户是否登陆
	userId, err := dao.QueryUserId(sessionId)
	//用户未登陆
	if err == sql.ErrNoRows {
		tool.InfoPrintln("user not login when request join or exit api", 0)
		return codes.NotLogin, nil
	}
	//数据库表操作失败
	if err != nil {
		tool.ErrorPrintln("sql query user id error", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	//在form表中插入或删除表项
	if action == 2 {
		err = dao.DeleteForm(userId, actId)
	} else {
		var form model.Form
		form.ActId, form.UserId = actId, userId
		form.JoinedAt = tool.GetTimeNowUnix()
		err = dao.InsertForm(form)
	}
	//操作数据库时出错
	if err != nil {
		tool.ErrorPrintln("sql when delete or insert into form table", userId, debug.Stack())
		return codes.MysqlError, nil
	}

	tool.InfoPrintln("user join or exit activity success", userId)
	return codes.OK, nil
}

//显示活动详情的一部分：显示参加某个活动的所有用户
func GetActivityUserList(actId uint) (codes.Code, []model.UserPublicMsg) {
	//活动id是必须的
	if actId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//从form表中获取所有加入活动的用户id
	rows, err := dao.QueryAllUsersId(actId)
	if err != nil {
		tool.ErrorPrintln("sql query user id error", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	//获取每个用户的信息
	var userId uint
	var objects []model.UserPublicMsg
	for rows.Next() {
		var obj model.UserPublicMsg
		err = rows.Scan(&userId)
		if err != nil {
			tool.ErrorPrintln("sql get user id error", 0, debug.Stack())
			return codes.MysqlError, nil
		}

		row := dao.QueryUsersMsg(userId)
		err := row.StructScan(&obj)
		if err != nil {
			tool.ErrorPrintln("sql query user msg error", 0, debug.Stack())
			return codes.MysqlError, nil
		}

		objects = append(objects, obj)
	}

	return codes.OK, objects
}

func GetCommentsList(actId uint, page uint) (codes.Code, []model.CommentListResponse) {
	//缺少活动的id参数
	if actId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//在comment_tab表中获取user_id、content、created_time
	rows, err := dao.QueryCommentMsg(actId, page)
	if err != nil {
		tool.ErrorPrintln("sql query comment msg", 0, debug.Stack())
		return codes.MysqlError, nil
	}

	var userId uint
	var objects []model.CommentListResponse
	for rows.Next() {
		var obj model.CommentListResponse
		rows.Scan(&userId, &obj.Content, &obj.CreatedAt)

		obj.Name, err = dao.QueryUserName(userId)
		if err != nil {
			tool.ErrorPrintln("sql query user name error", 0, debug.Stack())
			return codes.MysqlError, nil
		}

		objects = append(objects, obj)
	}

	return codes.OK, objects
}

func doActivityInfo(c *gin.Context) (codes.Code, interface{}) {
	var req model.ActivityInfoRequest
	err := c.BindJSON(&req)
	if err != nil {
		tool.ErrorPrintln("bind json error", 0, debug.Stack())
		return codes.BindJsonError, nil
	}
	sessionId, actId, page := req.SessionId, req.ActivityId, req.CommentPage

	//活动id是必要的参数
	if actId == 0 {
		tool.ErrorPrintln("request parameters is empty", 0, nil)
		return codes.MissParameter, nil
	}

	//如果用户是登陆状态，获取用户的userid
	var userId uint
	if sessionId != 0 {
		userId, err = dao.QueryUserId(sessionId)
		if err != nil {
			tool.ErrorPrintln("sql query user is login error", 0, debug.Stack())
			return codes.MysqlError, nil
		}
	}

	//在activities表中查询活动信息
	var obj model.ActivityInfoResponse
	var actDetail model.ActivityDetail
	row := dao.QueryActivityDetail(actId)
	err = row.StructScan(&actDetail)
	if err != nil {
		tool.ErrorPrintln("sql query activity deltail error", userId, debug.Stack())
		return codes.MysqlError, nil
	}
	obj.ActivityDetail = &actDetail

	//用户是否加入此活动
	if sessionId != 0 {
		obj.JoinStatus, err = dao.IsJoinin(userId, actId)
		if err != nil {
			tool.ErrorPrintln("sql query user is joinin", userId, debug.Stack())
			return codes.MysqlError, nil
		}
	} else {
		obj.JoinStatus = false
	}

	//查询加入activity id的所有用户
	var code codes.Code
	code, obj.UserList = GetActivityUserList(actId)
	if code != codes.OK {
		tool.ErrorPrintln("sql query user list error", userId, debug.Stack())
		return code, nil
	}

	//查询某个活动的评论列表
	code, obj.CommentList = GetCommentsList(actId, page)
	if code != codes.OK {
		tool.ErrorPrintln("sql query comment list error", userId, debug.Stack())
		return code, nil
	}

	return codes.OK, obj
}
