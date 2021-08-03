package handler

import (
	"database/sql"
	"example.com/greetings/codes"
	"example.com/greetings/constant"
	"example.com/greetings/dao"
	"example.com/greetings/model"
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
)



//用户注册接口
func Register(c *gin.Context) {
	//提取请求参数
	var req model.RegisterRequest
	c.BindJSON(&req)
	name, passwd, email, avatar := req.Name, req.Passwd, req.Email, req.Avatar

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	isMatch1, _ := regexp.MatchString(constant.NamePattern, name)
	isMatch2, _ := regexp.MatchString(constant.PasswdPattern, passwd)
	if name == "" || passwd == "" || email == "" || !isMatch1 || !isMatch2 {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//用户已经存在
	isRegister := dao.QueryUserIsRegister(name, email)
	if isRegister {
		res := ResponseFun(codes.UserExist, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.UserExist), res)
		return
	}

	//插入数据到user表中
	var user model.User
	user.Name, user.Avatar, user.Email = name, avatar, email
	user.Passwd = AddSalt(passwd)
	user.IsAdmin = false
	user.CreatedAt = GetTime()
	err := dao.InsertUser(user)

	//插入数据库时发生错误
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//注册成功
	res := ResponseFun(codes.OK, nil)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//用户登陆接口
func Login(c *gin.Context) {
	//提取请求参数
	var req model.LoginRequest
	c.BindJSON(&req)
	name, reqPasswd := req.Name, req.Passwd

	//传入的用户名和密码不能为空
	if name == "" || reqPasswd == "" {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	reqPasswd = AddSalt(reqPasswd) //密码加盐
	userId, passwd, err := dao.QueryUserIsExist(name)
	//用户不存在
	if err == sql.ErrNoRows {
		res := ResponseFun(codes.UserNotExist, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.UserNotExist), res)
		return
	}
	//密码错误
	if reqPasswd != passwd {
		res := ResponseFun(codes.PassWordError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.PassWordError), res)
		return
	}

	//插入session数据库
	var session model.Session
	session.Id = CreateSessionId(string(userId))
	session.UserId = userId
	err = dao.InsertSession(session)
	//数据库插入时出错
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}
	//登陆成功，返回sessoin
	res := ResponseFun(codes.OK, session)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//显示所有活动
func ShowActivities(c *gin.Context) {
	var req model.ShowActivityRequest
	c.BindJSON(&req)
	sessionId, page := req.SessionId, req.Page

	//如果用户是登陆状态，传入sessionId不为空
	var userId uint
	var err error
	if sessionId != ""{
		userId, err = dao.QueryUserId(sessionId)
		//数据库执行时错误
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}
		//参数错误，没有登陆
		if err == sql.ErrNoRows {
			res := ResponseFun(codes.NotLogin, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.NotLogin), res)
			return
		}
	}

	//获取活动表的当前page
	rows, err := dao.GetALLActivityRows(page)
	//查询时出错
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//构建返回数据
	var objects []model.ShowActivitiesRes
	for rows.Next() {
		var actId uint
		var act model.ShowActivitiesRes
		rows.Scan(&actId, &act.Title, &act.Start, &act.End)

		if req.SessionId != ""{
			act.JoinStatus, err = dao.IsJoinin(userId, actId)
			if err != nil {
				res := ResponseFun(codes.MysqlError, nil)
				c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
				return
			}
		}else{
			act.JoinStatus = false
		}

		objects = append(objects, act)
	}

	//处理完成，返回活动列表objects
	res := ResponseFun(codes.OK, objects)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//活动过滤器
func ActivitiesSelector(c *gin.Context) {
	var req model.ActivitySelectorRequest
	c.BindJSON(&req)
	sessionId, actType, start, end, page := req.SessionId, req.Type, req.Start, req.End, req.Page

	//如果用户是登陆状态，获取user_id
	var userId uint
	var err error
	if sessionId != ""{
		userId, err = dao.QueryUserId(sessionId)
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}
	}

	//封装sql语句，得到符合条件的活动的行
	rows, err := dao.SqlActivitiesSelector(actType, start, end, page)
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//遍历等到activity_selector的id、type_id、title、start_time、end_time
	var actId uint
	var objects []model.ActivitySelectorResponse
	for rows.Next() {
		var act model.ActivitySelectorResponse
		rows.Scan(&actId, &act.Title, &act.Start, &act.End)

		//是否参加活动
		if sessionId != "" {
			uintSessionId, _ := strconv.Atoi(sessionId)
			act.JoinStatus, err = dao.IsJoinin(userId, uint(uintSessionId))
			if err != nil {
				res := ResponseFun(codes.MysqlError, nil)
				c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
				return
			}
		}else{
			act.JoinStatus = false
		}

		objects = append(objects, act)
	}
	res := ResponseFun(codes.OK, objects)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//发表评论
func CreateComment(c *gin.Context) {
	var req model.CreateCommentRequest
	c.BindJSON(&req)
	sessionId, activityId, content := req.SessionId, req.ActivityId, req.Content

	//传入内容不能为空
	if sessionId == "" || string(activityId) == "" || content == "" {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//从session表中获取用户user id，查询不到session id对应的用户时也会报错
	userId, err := dao.QueryUserId(sessionId)
	//数据库查表时出错
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//将数据插入comment表中
	var comment model.Comment
	comment.UserId, comment.ActivityId, comment.Content = userId, activityId, content
	comment.CreatedAt = GetTime()
	err = dao.InsertComment(comment)
	//插入数据库表时出错
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	res := ResponseFun(codes.OK, nil)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//显示活动的所有评论
func CommentsList(c *gin.Context) {
	var req model.CommentListRequest
	c.BindJSON(&req)
	actId, page := req.ActivityId, req.Page

	//缺少活动的id参数
	if actId == "" {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//在comment_tab表中获取user_id、content、created_time
	rows, err := dao.QueryCommentMsg(actId, page)
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	var userId uint
	var objects []model.CommentListResponse
	for rows.Next() {
		var obj model.CommentListResponse
		rows.Scan(&userId, &obj.Content, &obj.CreatedAt)

		obj.Name, err = dao.QueryUserName(userId)
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}

		objects = append(objects, obj)
	}

	res := ResponseFun(codes.OK, objects)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}

//显示加入某个活动的所有用户
func ActivityUserList(c *gin.Context) {
	var req model.ActivityUserListRequest
	c.BindJSON(&req)
	actId := req.ActivityId

	//不可缺少活动id
	if actId == "" {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//从form表中获取所有加入活动的用户id
	rows, err := dao.QueryAllUsersId(actId)
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//获取每个用户的信息
	var userId uint
	var objects []model.ActivityUserListResponse
	for rows.Next() {
		var obj model.ActivityUserListResponse
		err = rows.Scan(&userId)
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}

		err = dao.QueryUsersMsg(userId, &obj)
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}

		objects = append(objects, obj)
	}

	res := ResponseFun(codes.OK, objects)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//显示用户加入的所有活动
func ShowJoinedActivities(c *gin.Context) {
	var req model.ShowJoinedActivitiesRequest
	c.BindJSON(&req)
	sessionId := req.SessionId

	//用户需要在登陆状态，session_id不能为空
	if sessionId == "" {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//获取用户user_id
	userId, err := dao.QueryUserId(sessionId)
	if err == sql.ErrNoRows {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//获取用户加入的所有活动的id
	rows, err := dao.GetAllJoinActivities(userId)
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//获取活动的信息，标题、开始时间、结束时间
	var actId uint
	var objects []model.ShowJoinedActivitiesResponse
	for rows.Next() {
		var obj model.ShowJoinedActivitiesResponse
		rows.Scan(&actId)

		err := dao.QueryActivityMsg(actId, &obj)
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}

		objects = append(objects, obj)
	}

	res := ResponseFun(codes.OK, objects)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}


//用户加入或退出活动
func JoinOrExit(c *gin.Context) {
	var req model.JoinOrExitRequest
	c.BindJSON(&req)
	sessionId, actId, actionStr := req.SessionId, req.ActivityId, req.Action

	//传入参数都不能为空或格式错误
	action, err := strconv.Atoi(actionStr)
	if sessionId == "" || actId == "" || actionStr == "" || err != nil{
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//查看用户是否登陆
	userId, err := dao.QueryUserId(sessionId)
	//用户未登陆
	if err == sql.ErrNoRows {
		res := ResponseFun(codes.NotLogin, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.NotLogin), res)
		return
	}
	//数据库表操作失败
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//在form表中插入或删除表项
	if action == 0 {
		err = dao.DeleteForm(userId, actId)
	} else {
		var form model.Form
		form.ActId, form.UserId = actId, userId
		form.JoinedAt = GetTime()
		err = dao.AddForm(form)
	}
	//操作数据库时出错
	if err != nil{
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	res := ResponseFun(codes.OK, nil)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
	return
}

/*
func ActivityInfo(c *gin.Context) {
	var req model.ActivityInfoRequest
	c.BindJSON(&req)
	sessionId, actId := req.SessionId, req.ActivityId

	//活动id是必要的参数
	if actId == 0 {
		res := ResponseFun(codes.MissParameter, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter), res)
		return
	}

	//如果用户是登陆状态，获取用户的userid
	var userId uint
	var err error
	if sessionId != 0 {
		userId, err = dao.QueryUserId(string(sessionId))
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}
	}

	//在activities表中查询活动信息
	var obj model.ActivityInfoResponse
	err = dao.QueryActivityDetail(actId, &obj)
	if err != nil {
		res := ResponseFun(codes.MysqlError, nil)
		c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
		return
	}

	//用户是否加入此活动
	if sessionId != 0 {
		obj.JoinStatus, err = dao.IsJoinin(userId, actId)
		if err != nil {
			res := ResponseFun(codes.MysqlError, nil)
			c.JSON(codes.HTTPStatusFromCode(codes.MysqlError), res)
			return
		}
	} else {
		obj.JoinStatus = false
	}

	//在form表中查询所有加入activity id的用户
	str = "select userId from form_tab where actId='" + actId + "'"
	rows, _ := globalVariable.DB.Query(str)

	for rows.Next() {
		var user main2.UserMsg
		var uid uint
		rows.Scan(&uid)

		//获取每个参加活动的用户的名称和头像，
		str = "select name, image from user_tab where id='" + userId + "'"
		globalVariable.DB.QueryRow(str).Scan(user.Name, user.Image)

		//将user_msg对象append进activity_detail的用户列表中
		ad.Users = append(ad.Users, user)
	}
	c.JSON(http.StatusOK, ad)
}
*/
