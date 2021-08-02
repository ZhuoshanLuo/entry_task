package handler

import (
	"database/sql"
	"example.com/greetings/codes"
	"example.com/greetings/constant"
	"example.com/greetings/dao"
	main2 "example.com/greetings/dir1"
	"example.com/greetings/globalVariable"
	"example.com/greetings/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func Activity_info(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)

	activity_id := json["id"].(string)
	user_id := json["user_id"].(string)

	//活动id不能为空
	if activity_id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "activity id cannot be empty!",
		})
		return
	}

	//在activities表中查询activity id
	var ad main2.Activity_detail
	str := "select title, content, location, start_time, end_time from activities_tab where id='" + activity_id + "'"
	globalVariable.DB.QueryRow(str).Scan(&ad.Title, &ad.Content, &ad.Location, &ad.Start_time, &ad.End_time)

	//用户是否加入此活动
	if user_id != "" {
		str = "select * from form_tab where user_id='" + user_id + "' and id='" + activity_id + "'"
		row, _ := globalVariable.DB.Query(str)
		if row.Next() == false {
			ad.Join_status = 0
		} else {
			ad.Join_status = 1
		}
	} else {
		ad.Join_status = 0
	}

	//在form表中查询所有加入activity id的用户
	str = "select user_id from form_tab where activity_id='" + activity_id + "'"
	rows, _ := globalVariable.DB.Query(str)

	for rows.Next() {
		var user main2.UserMsg
		var uid uint
		rows.Scan(&uid)

		//获取每个参加活动的用户的名称和头像，
		str = "select name, image from user_tab where id='" + user_id + "'"
		globalVariable.DB.QueryRow(str).Scan(user.Name, user.Image)

		//将user_msg对象append进activity_detail的用户列表中
		ad.Users = append(ad.Users, user)
	}
	c.JSON(http.StatusOK, ad)
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("CheckErr:%v", err)
	}
}

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
		c.JSON(codes.HTTPStatusFromCode(codes.MissParameter,), res)
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
	res := ResponseFun(codes.OK, session.Id)
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
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
}

func Joined_activities_view(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	session_id := json["session_id"].(string)

	//用户需要在登陆状态，session_id不能为空
	if session_id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "User not login!",
		})
		return
	}

	//在session_tab中获取用户user_id
	var user_id uint
	str := "select user_id from session_tab where session_id=" + session_id
	err := globalVariable.DB.QueryRow(str).Scan(&user_id)
	CheckErr(err)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "user not login",
		})
		return
	}

	//在form_tab中获取用户加入的所有活动的id
	str = "select id from form_tab where user_id=" + string(user_id)
	rows, _ := globalVariable.DB.Query(str)

	//在activities_tab中获取活动的信息，标题、开始时间、结束时间
	var act_id uint
	var objects []main2.Activities_joinin
	for rows.Next() {
		var obj main2.Activities_joinin
		rows.Scan(&act_id)

		str = "select title, start_time, end_time from activities_tab where id=" + string(act_id)
		globalVariable.DB.QueryRow(str).Scan(&obj.Title, &obj.Start_time, &obj.End_time)

		objects = append(objects, obj)
	}

	c.JSON(http.StatusOK, objects)
}

func Join_quit(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)

	session_id := json["session_id"].(string)
	activity_id := json["activity_id"].(string)
	action := json["action"].(uint)

	//查看用户是否登陆
	var user_id uint
	err := globalVariable.DB.QueryRow("select user_id from session_tab where session_id='" + session_id + "'").Scan(&user_id)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "user not login!",
		})
		return
	}

	//在form表中插入或删除表项
	if action == 0 {
		globalVariable.DB.Exec("delete from form_tab where user_id='" + string(user_id) + "' and activity_id='" + activity_id + "'")
	} else {
		globalVariable.DB.Exec("insert into form_tab(id, activity_id, user_id, join_at) values(1, '" + activity_id + "', '" + string(user_id) + "', '" + string(time.Now().Unix()) + "'")
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
	})
}

func Manage_login(c *gin.Context) {

}

func Manage_add_activity(c *gin.Context) {

}

func Manage_del_activity(c *gin.Context) {

}

func Manage_add_activity_type(c *gin.Context) {

}

func Manage_edit_activity(c *gin.Context) {

}

func Manage_show_users(c *gin.Context) {

}

func Manage_edit_activity_type(c *gin.Context) {

}

func Manage_del_activity_type(c *gin.Context) {

}
