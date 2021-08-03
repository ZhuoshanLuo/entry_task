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

func InsertUser(user model.User) error {
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

func InsertSession(session model.Session) error {
	str := fmt.Sprintf("insert into session_tab(id, user_id) values(%d, %d)", session.Id, session.UserId)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func GetALLActivityRows(page uint) (*sql.Rows, error) {
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

func InsertComment(obj model.Comment) error {
	str := fmt.Sprintf("insert into comments_tab(user_id, activity_id, content, created_at) values(%d, %d, '%s', %d)", obj.UserId, obj.ActivityId, obj.Content, obj.CreatedAt)
	_, err := globalVariable.DB.Exec(str)
	if err != nil {
		return err
	}
	return nil
}

func QueryUserId(sessionId uint) (uint, error) {
	var userId uint
	str := fmt.Sprintf("select user_id from session_tab where id=%d", sessionId)
	err := globalVariable.DB.QueryRow(str).Scan(&userId)
	return userId, err
}

func SqlActivitiesSelector(actType string, start, end, page uint) (*sql.Rows, error) {
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
	if start != 0 {
		if flag == 0 {
			flag = 1
			str += "where a.start_time>=" + string(start) + " "
		} else {
			str += "and a.start_time>=" + string(start) + " "
		}
	}
	if end != 0 {
		if flag == 0 {
			flag = 1
			str += "where a.end_time<=" + string(end) + " "
		} else {
			str += "and a.end_time<=" + string(end) + " "
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

func AddForm(form model.Form) error {
	str := fmt.Sprintf("insert into form_tab(activity_id, user_id, joined_at) values(%d, %d, %d)", form.ActId, form.UserId, form.JoinedAt)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func DeleteForm(userId uint, actId uint) error {
	str := fmt.Sprintf("delete from form_tab where user_id=%d and activity_id=%d", userId, actId)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func QueryActivityMsg(actId uint) (string, uint, uint, error) {
	var title string
	var start, end uint
	str := fmt.Sprintf("select title, start_time, end_time from activities_tab where id=%d", actId)
	err := globalVariable.DB.QueryRow(str).Scan(&title, &start, &end)
	return title, start, end, err
}

func GetAllJoinActivities(userId uint) (*sql.Rows, error) {
	str := fmt.Sprintf("select activity_id from form_tab where user_id=%d", userId)
	rows, err := globalVariable.DB.Query(str)
	return rows, err
}

func QueryUsersMsg(userId uint, obj *model.ActivityUserListResponse) error {
	str := fmt.Sprintf("select name, avatar from user_tab where id=%d", userId)
	err := globalVariable.DB.QueryRow(str).Scan(&obj.Name, &obj.Avatar)
	return err
}

func QueryAllUsersId(actId string) (*sql.Rows, error) {
	str := fmt.Sprintf("select user_id from form_tab where activity_id=%s", actId)
	return globalVariable.DB.Query(str)
}

func QueryActivityDetail(actId uint) (*model.ActivityDetail, error) {
	var data model.ActivityDetail
	str := fmt.Sprintf("select title, content, location, start_time, end_time from activities_tab where id=%d", actId)
	err := globalVariable.DB.QueryRow(str).Scan(&data.Title, &data.Start, &data.End, &data.Location, &data.Location, &data.Content)
	return &data, err
}

func InsertActvity(act model.Activity) error {
	str := fmt.Sprintf("insert into activities_tab(type_id, title, content, locatoin, start_time, end_time) values(%d, '%s', '%s', '%s', '%d', '%d')", act.TypeId, act.Title, act.Content, act.Location, act.Start, act.End)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func QueryUserAdmin(userId uint) (bool, error) {
	var isAdmin bool
	str := fmt.Sprintf("select is_admin from user_tab where id=%d", userId)
	err := globalVariable.DB.QueryRow(str).Scan(&isAdmin)
	return isAdmin, err
}

func DeleteActivity(actId uint) error {
	str := fmt.Sprintf("delete from activities_tab where id=%d", actId)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func UpdateActivity(id, typeId uint, title, content, location string, start, end uint) error {
	var flag = 0
	str := "update activities_tab "
	if typeId != 0 {
		if flag == 0 {
			str += "set type_id=" + string(typeId) + " "
			flag = 1
		} else {
			str += ", type_id=" + string(typeId) + " "
		}
	}
	if title != "" {
		if flag == 0 {
			str += "set title=" + string(title) + " "
			flag = 1
		} else {
			str += ", title=" + string(title) + " "
		}
	}
	if content != "" {
		if flag == 0 {
			str += "set content=" + string(content) + " "
			flag = 1
		} else {
			str += ", content=" + string(content) + " "
		}
	}
	if location != "" {
		if flag == 0 {
			str += "set location=" + string(location) + " "
			flag = 1
		} else {
			str += ", location=" + string(location) + " "
		}
	}
	if start != 0 {
		if flag == 0 {
			str += "set start_time=" + string(start) + " "
			flag = 1
		} else {
			str += ", start_time=" + string(start) + " "
		}
	}
	if end != 0 {
		if flag == 0 {
			str += "set end_time=" + string(end) + " "
			flag = 1
		} else {
			str += ", end_time=" + string(end) + " "
		}
	}
	str += "where id=" + string(id)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func InsertActivityType(obj model.ActivityType) error {
	str := fmt.Sprintf("insert into activities_type_tab(name) values('%s')", obj.Name)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func DeleteActivityType(typeId uint) error {
	str := fmt.Sprintf("delete from activities_type_tab where id=%d", typeId)
	_, err := globalVariable.DB.Exec(str)
	return err
}

func QueryAllActivityType() (*sql.Rows, error) {
	str := "select name from activities_type_tab"
	return globalVariable.DB.Query(str)
}

func QueryAllUserMsg() (*sql.Rows, error) {
	str := "select name, email, avatar from user_tab"
	return globalVariable.DB.Query(str)
}

func UpdateActivityType(id uint, name string) error {
	str := "update activities_type_tab "
	if name != "" {
		str += "set name=" + name
	}
	str += "where id=" + string(id)
	_, err := globalVariable.DB.Exec(str)
	return err
}
