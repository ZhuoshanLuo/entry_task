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
	str := fmt.Sprintf("insert into user_tab(id, name, passwd, email, Image, is_admin, created_at) values ('%s', '%s', '%s', '%s', %t, %d)", user.Name, user.Passwd, user.Email, user.Avatar, user.IsAdmin, user.CreatedAt)
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
	str := fmt.Sprintf("insert into comment_tab(id, user_id, activityId, content, created_at) values(%d, %d, '%s', %d", obj.Id, obj.ActivityId, obj.Content, obj.CreatedAt)
	_, err := globalVariable.DB.Exec(str)
	if err != nil {
		return err
	}
	return nil
}

func QueryUserId(sessionId string) (uint, error){
	var userId uint
	str := fmt.Sprintf("select user_id from session_tab where session_id=%s", sessionId)
	err := globalVariable.DB.QueryRow(str).Scan(&userId)
	return userId, err
}
