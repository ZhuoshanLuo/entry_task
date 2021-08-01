package handler

import (
	"database/sql"
	"example.com/greetings/codes"
	"example.com/greetings/dao"
	main2 "example.com/greetings/dir1"
	"example.com/greetings/globalVariable"
	"example.com/greetings/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

func Query_activity_selector(act_type, start, end string) string {
	return ""
}

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("CheckErr:%v", err)
	}
}

//用户注册接口
func Register(c *gin.Context) {
	//提取请求参数
	var regReq model.RegisterRequest
	c.BindJSON(&regReq)
	name, passwd, email, avatar := regReq.Name, regReq.Passwd, regReq.Email, regReq.Avatar

	//传入的用户名、密码和邮箱不能为空或不符合pattern（头像可为空）
	//isMatch1, _ := regexp.MatchString(constant.NamePattern, name)
	//isMatch2, _ := regexp.MatchString(constant.PasswdPattern, passwd)
	//|| !isMatch1 || !isMatch2
	if name == "" || passwd == "" || email == "" {
		res := model.Response{
			Code: codes.Fail,
			Msg:  codes.Errorf(codes.Fail),
		}
		c.JSON(codes.HTTPStatusFromCode(codes.Fail), res)
		return
	}

	//user表中用户名和邮箱不能存在
	isExist := dao.QueryUserIsExist(name, email)
	if isExist {
		res := model.Response{
			Code: codes.Fail,
			Msg:  codes.Errorf(codes.Fail),
		}
		c.JSON(codes.HTTPStatusFromCode(codes.Fail), res)
		return
	}

	//插入数据到user表中
	var user model.User
	user.Name, user.Avatar, user.Email = name, avatar, email
	user.Id = CreateId(name)
	user.Passwd = AddSalt(passwd)
	user.IsAdmin = false
	user.CreatedAt = GetTime()

	dao.InsertUser(user)
	res := model.Response{
		Code: codes.OK,
		Msg:  codes.Errorf(codes.OK),
	}
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
}

//用户登陆接口
func Login(c *gin.Context) {
	//提取请求参数
	var loginReq model.LoginRequest
	c.BindJSON(&loginReq)
	name, passwd := loginReq.Name, loginReq.Passwd

	//传入的用户名和密码不能为空
	if name == "" || passwd == "" {
		res := model.Response{
			Code: codes.Fail,
			Msg:  codes.Errorf(codes.Fail),
		}
		c.JSON(codes.HTTPStatusFromCode(codes.Fail), res)
		return
	}

	//查看mysql，用户名或密码错误，登陆失败
	passwd = AddSalt(passwd) //密码加盐
	userId, err := dao.QueryIdFromUsertabWithName(name)
	if err == sql.ErrNoRows {
		res := model.Response{
			Code: codes.Fail,
			Msg:  codes.Errorf(codes.Fail),
		}
		c.JSON(codes.HTTPStatusFromCode(codes.Fail), res)
		return
	}

	//登陆成功，保存并返回session id
	var session model.Session
	session.SessionId = CreateId(string(userId))
	dao.InsertSession(session.SessionId, userId)
	res := model.Response{
		Code: codes.OK,
		Msg:  codes.Errorf(codes.OK),
		Data: session,
	}
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
}

func Show_activities(c *gin.Context) {
	var req model.ShowActivityRequest
	c.BindJSON(&req)
	user_id := req.UserId

	//获取活动表的所有行
	rows := dao.GetALLActivityRows()

	var objects []model.ShowActivitiesRes
	for rows.Next() {
		var actId uint
		var act model.ShowActivitiesRes
		rows.Scan(&actId, &act.Title, &act.Start, &act.End)
		act.JoinStatus = dao.IsJoinin(user_id, actId)

		objects = append(objects, act)
	}

	res := model.Response{
		Code: codes.OK,
		Msg:  codes.Errorf(codes.OK),
		Data: objects,
	}
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
}

func Activities_selector(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	req_type_id := json["type_id"].(string)
	req_start_time := json["start_time"].(string)
	req_end_time := json["end_time"].(string)
	user_id := json["user_id"].(string)

	//封装sql语句，得到符合条件的活动
	//返回包括type_id、title、start_time、end_time
	str := Query_activity_selector(req_type_id, req_start_time, req_end_time)
	rows, _ := globalVariable.DB.Query(str)

	//遍历等到activity_selector的id、type_id、title、start_time、end_time
	var act_id uint
	var type_id uint
	var objects []main2.Activity_select
	for rows.Next() {
		var as main2.Activity_select
		rows.Scan(&act_id, &type_id, &as.Title, &as.Start_time, &as.End_time)

		//根据type_id查询活动类型
		str = "select name from activities_type_tab where id=" + string(type_id)
		globalVariable.DB.QueryRow(str).Scan(&as.Type_name)

		//根据user_id查询form_tab是否参加活动
		str = "select * from form_tab where id=" + string(act_id) + " and user_id=" + user_id
		_, err := globalVariable.DB.Query(str)
		if err == sql.ErrNoRows {
			as.Join_status = 0
		} else {
			as.Join_status = 1
		}

		objects = append(objects, as)
	}
	c.JSON(http.StatusOK, objects)
}

func Create_comment(c *gin.Context) {
	var req model.CreateCommentRequest
	c.BindJSON(&req)
	sessionId, activityId, content := req.SessionId, req.ActivityId, req.Content

	//传入内容不能为空
	if string(sessionId) == "" || string(activityId) == "" || content == "" {
		res := model.Response{
			Code: codes.Fail,
			Msg:  codes.Errorf(codes.Fail),
		}
		c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
		return
	}

	//从session表中获取用户user id
	userId := dao.GetUserIdFromSession(sessionId)

	//将数据插入comment表中
	var comment model.Comment
	comment.Id = CreateId("")
	comment.CreatedAt = GetTime()
	comment.UserId, comment.ActivityId, comment.Content = userId, activityId, content
	err := dao.InsertComment(comment)
	if err != nil {
		res := model.Response{
			Code: codes.Fail,
			Msg:  codes.Errorf(codes.Fail),
		}
		c.JSON(codes.HTTPStatusFromCode(codes.Fail), res)
		return
	}
	res := model.Response{
		Code: codes.OK,
		Msg:  codes.Errorf(codes.OK),
	}
	c.JSON(codes.HTTPStatusFromCode(codes.OK), res)
}

//显示活动的所有评论，暂不考虑多级评论和评论分页
func Comments(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)
	activity_id := json["activity_id"].(string)

	//活动的id是必须的
	if activity_id == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "activity id is null!",
		})
		return
	}

	//在comment_tab表中获取user_id、content、created_time
	//遍历查询结果，根据user_id在user_tab中得到用户name
	str := "select user_id, content, created_at from comments_tab where activity_id=" + activity_id
	rows, _ := globalVariable.DB.Query(str)

	var user_id uint
	var objects []main2.Comment
	for rows.Next() {
		var obj main2.Comment
		rows.Scan(&user_id, &obj.Content, &obj.Created_time)

		str = "select name from user_tab where id=" + string(user_id)
		globalVariable.DB.QueryRow(str).Scan(obj.User_name)

		objects = append(objects, obj)
	}

	c.JSON(http.StatusOK, objects)
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
