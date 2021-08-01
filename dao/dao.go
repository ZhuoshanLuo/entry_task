package dao

import (
	"database/sql"
	"example.com/greetings/globalVariable"
	"example.com/greetings/model"
	"fmt"
)

func QueryUserIsExist(name string, email string) bool {
	str := fmt.Sprintf("select * from user_tab where name='%s' or email='%s'", name, email)
	row, _ := globalVariable.DB.Query(str)
	if row.Next() {
		return true
	}
	return false
}

func InsertUser(user model.User) {
	str := fmt.Sprintf("insert into user_tab(id, name, passwd, email, Image, is_admin, created_at) values (%d, '%s', '%s', '%s', '%s', %t, %d)", user.Id, user.Name, user.Passwd, user.Email, user.Avatar, user.IsAdmin, user.CreatedAt)
	_, err := globalVariable.DB.Exec(str)
	if err != nil {
		return 
	}
}

func QueryIdFromUsertabWithName(name string) (uint, error) {
	var userId uint
	str := fmt.Sprintf("select id from user_tab where name='%s'", name)
	err := globalVariable.DB.QueryRow(str).Scan(&userId)
	return userId, err
}

func InsertSession(sessionId uint, userId uint) {
	str := fmt.Sprintf("insert into session_tab(id, user_id) values(%d, %d)", sessionId, userId)
	globalVariable.DB.Exec(str)
}

func GetALLActivityRows() *sql.Rows {
	str := fmt.Sprintf("select id, title, start_time, end_time from activities_tab")
	rows, _ := globalVariable.DB.Query(str)
	return rows
}

func IsJoinin(userId string, actId uint) uint {
	str := fmt.Sprintf("select * from form_tab where activity_id=%d and user_id=%d", actId, userId)
	if userId != "" {
		row, _ := globalVariable.DB.Query(str)
		if row.Next() == false {
			return 0
		} else {
			return 1
		}
	}
	return 0
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
