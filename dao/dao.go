package dao

import (
	"database/sql"
	"example.com/greetings/constant"
	"example.com/greetings/globalVariable"
	"example.com/greetings/model"
	"fmt"
)

func QueryUserIsRegister(name string, email string) bool {
	str := fmt.Sprintf("select * from user_tab where name='%s' or email='%s'", name, email)
	row, _ := globalVariable.DB.Query(str)
	if row.Next() {
		return true
	}
	return false
}

func InsertUser(user model.User) error{
	str := fmt.Sprintf("insert into user_tab(name, passwd, email, avatar, is_admin, created_at) values ('%s', '%s', '%s', '%s', %t, %d)", user.Name, user.Passwd, user.Email, user.Avatar, user.IsAdmin, user.CreatedAt)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func QueryUserIsExist(name string) (uint, string, error) {
	var userId uint
	var passwd string
	str := fmt.Sprintf("select id, passwd from user_tab where name='%s'", name)
	err := globalVariable.DB.QueryRow(str).Scan(&userId, &passwd)
	return userId, passwd, err
}

func InsertSession(session model.Session) error{
	str := fmt.Sprintf("insert into session_tab(id, user_id) values(%d, %d)", session.Id, session.UserId)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func GetALLActivityRows(page uint) (*sql.Rows,error) {
	offset := page * 10
	str := fmt.Sprintf("select id, title, start_time, end_time from activities_tab limit %d offset %d", constant.Limit, offset)
	rows, err := globalVariable.DB.Query(str)
	return rows, err
}

func IsJoinin(userId uint, actId uint) (bool, error) {
	str := fmt.Sprintf("select * from form_tab where activity_id=%d and user_id=%d", actId, userId)
	row, err := globalVariable.DB.Query(str)
	if row.Next() {
		return true, err
	}
	return false, err
}

func GetUserIdFromSession(sessionId uint) uint {
	var userId uint
	str := fmt.Sprintf("select user_id from session_tab where session_id=%d", sessionId)
	globalVariable.DB.QueryRow(str).Scan(&userId)
	return userId
}

func InsertComment(obj model.Comment) error {
	str := fmt.Sprintf("insert into comments_tab(user_id, activity_id, content, created_at) values(%d, %d, '%s', %d)", obj.UserId, obj.ActivityId, obj.Content, obj.CreatedAt)
	_, err := globalVariable.DB.Exec(str)
	if err != nil {
		return err
	}
	return nil
}

func QueryUserId(sessionId string) (uint, error){
	var userId uint
	str := fmt.Sprintf("select user_id from session_tab where id=%s", sessionId)
	err := globalVariable.DB.QueryRow(str).Scan(&userId)
	return userId, err
}

func SqlActivitiesSelector(actType, start, end string, page uint) (*sql.Rows, error) {
	var flag = 0
	offset := page * 10
	str := "select a.id, a.title, a.start_time, a.end_time from activities_tab a left join activities_type_tab b on a.type_id=b.id "
	if actType != "" {
		if flag == 0 {
			flag = 1
			str += "where b.name='" + actType + "' "
		} else {
			str += "and b.name='" + actType + "' "
		}
	}
	if start != "" {
		if flag == 0 {
			flag = 1
			str += "where a.start_time>=" + start + " "
		} else {
			str += "and a.start_time>=" + start + " "
		}
	}
	if end != "" {
		if flag == 0 {
			flag = 1
			str += "where a.end_time<=" + end + " "
		} else {
			str += "and a.end_time<=" + end + " "
		}
	}
	s := fmt.Sprintf("limit 10 offset %d", offset)
	str += s
	rows, err := globalVariable.DB.Query(str)
	return rows, err
}

func QueryCommentMsg(actId string, page uint) (*sql.Rows, error) {
	offset := page * 10
	str := fmt.Sprintf("select user_id, content, created_at from comments_tab where activity_id=%s limit 10 offset %d", actId, offset)
	rows, err := globalVariable.DB.Query(str)
	return rows, err
}

func QueryUserName(userId uint) (string, error) {
	var userName string
	str := fmt.Sprintf("select name from user_tab where id=%d", userId)
	err := globalVariable.DB.QueryRow(str).Scan(&userName)
	return userName, err
}
