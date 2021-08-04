package dao

import (
	"database/sql"
	"example.com/greetings/globalVariable"
	"example.com/greetings/model"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func QueryUserIsRegister(name string, email string) (bool, error) {
	sql := "select * from user_tab where name=? or email=?"
	row, err := globalVariable.DB.Queryx(sql, name, email)
	//数据库发生错误
	if err != nil || !row.Next() {
		return false, err
	}
	return true, err
}

func InsertUser(user model.User) error {
	sql := "insert into user_tab(name, passwd, email, avatar, is_admin, created_at) values(:name, :passwd, :email, :avatar, :is_admin, :created_at)"
	_, err := globalVariable.DB.NamedExec(sql, &user)
	return err
}

func QueryUserIsExist(name string) (uint, string, error) {
	var userId uint
	var passwd string
	sql := "select id, passwd from user_tab where name=?"
	row := globalVariable.DB.QueryRowx(sql, name)
	err := row.Scan(&userId, &passwd)
	return userId, passwd, err
}

func InsertSession(session model.Session) error {
	sql := "insert into session_tab(id, user_id) values(:id, :user_id)"
	_, err := globalVariable.DB.NamedExec(sql, &session)
	return err
}

func GetALLActivityRows(page uint) (*sqlx.Rows, error) {
	offset := page * 10
	sql := "select id, title, start_time, end_time from activities_tab limit ? offset ?"
	rows, err := globalVariable.DB.Queryx(sql, 10, offset)
	return rows, err
}

func IsJoinin(userId uint, actId uint) (bool, error) {
	sql := "select * from form_tab where activity_id=? and user_id=?"
	row, err := globalVariable.DB.Queryx(sql, actId, userId)
	if row.Next() {
		return true, err
	}
	return false, err
}

func InsertComment(obj model.Comment) error {
	sql := "insert into comments_tab(user_id, activity_id, content, created_at) values(:user_id, :activity_id, :content, :created_at)"
	_, err := globalVariable.DB.NamedExec(sql, &obj)
	if err != nil {
		return err
	}
	return nil
}

func QueryUserId(sessionId uint) (uint, error) {
	var userId uint
	sql := "select user_id from session_tab where id=?"
	err := globalVariable.DB.Get(&userId, sql, sessionId)
	return userId, err
}

func SqlActivitiesSelector(actType string, start, end, page uint) (*sqlx.Rows, error) {
	var flag = 0
	offset := page * 10
	sql := "select a.id, a.title, a.start_time, a.end_time from activities_tab a left join activities_type_tab b on a.type_id=b.id "
	if actType != "" {
		if flag == 0 {
			flag = 1
			sql += "where b.name='" + actType + "' "
		} else {
			sql += "and b.name='" + actType + "' "
		}
	}
	if start != 0 {
		if flag == 0 {
			flag = 1
			sql += "where a.start_time>=" + string(start) + " "
		} else {
			sql += "and a.start_time>=" + string(start) + " "
		}
	}
	if end != 0 {
		if flag == 0 {
			flag = 1
			sql += "where a.end_time<=" + string(end) + " "
		} else {
			sql += "and a.end_time<=" + string(end) + " "
		}
	}
	s := fmt.Sprintf("limit 10 offset %d", offset)
	sql += s
	rows, err := globalVariable.DB.Queryx(sql)
	return rows, err
}

func QueryCommentMsg(actId uint, page uint) (*sqlx.Rows, error) {
	offset := page * 10
	sql := "select user_id, content, created_at from comments_tab where activity_id=? limit ? offset ?"
	rows, err := globalVariable.DB.Queryx(sql, actId, 10, offset)
	return rows, err
}

func QueryUserName(userId uint) (string, error) {
	var userName string
	sql := "select name from user_tab where id=?"
	err := globalVariable.DB.QueryRowx(sql, userId).Scan(&userName)
	return userName, err
}

func AddForm(form model.Form) error {
	sql := "insert into form_tab(activity_id, user_id, joined_at) values(:activity_id, :user_id, :joined_at)"
	_, err := globalVariable.DB.NamedExec(sql, &form)
	return err
}

func DeleteForm(userId uint, actId uint) error {
	sql := "delete from form_tab where user_id=? and activity_id=?"
	_, err := globalVariable.DB.Exec(sql, userId, actId)
	return err
}

func QueryActivityMsg(actId uint) (string, uint, uint, error) {
	var title string
	var start, end uint
	sql := "select title, start_time, end_time from activities_tab where id=?"
	row := globalVariable.DB.QueryRowx(sql, actId)
	err := row.Scan(&title, &start, &end)
	return title, start, end, err
}

func GetAllJoinActivities(userId uint) (*sqlx.Rows, error) {
	sql := "select activity_id from form_tab where user_id=?"
	rows, err := globalVariable.DB.Queryx(sql, userId)
	return rows, err
}

func QueryUsersMsg(userId uint) *sqlx.Row {
	sql := "select name, avatar from user_tab where id=?"
	row := globalVariable.DB.QueryRowx(sql, userId)
	return row
}

func QueryAllUsersId(actId uint) (*sqlx.Rows, error) {
	sql := "select user_id from form_tab where activity_id=?"
	return globalVariable.DB.Queryx(sql, actId)
}

func QueryActivityDetail(actId uint) (*sqlx.Row) {
	sql := "select title, content, location, start_time, end_time from activities_tab where id=?"
	return globalVariable.DB.QueryRowx(sql, actId)
}

func InsertActvity(act model.Activity) error {
	str := fmt.Sprintf("insert into activities_tab(type_id, title, content, location, start_time, end_time) values(%d, '%s', '%s', '%s', '%d', '%d')", act.TypeId, act.Title, act.Content, act.Location, act.Start, act.End)
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
	sid := fmt.Sprintf("%d", id)
	stypeId := fmt.Sprintf("%d", typeId)
	sstart := fmt.Sprintf("%d", start)
	send := fmt.Sprintf("%d", end)
	str := "update activities_tab "
	if typeId != 0 {
		if flag == 0 {
			str += "set type_id=" + stypeId + " "
			flag = 1
		} else {
			str += ", type_id=" + stypeId + " "
		}
	}
	if title != "" {
		if flag == 0 {
			str += "set title='" + string(title) + "' "
			flag = 1
		} else {
			str += ", title='" + string(title) + "' "
		}
	}
	if content != "" {
		if flag == 0 {
			str += "set content='" + string(content) + "' "
			flag = 1
		} else {
			str += ", content='" + string(content) + "' "
		}
	}
	if location != "" {
		if flag == 0 {
			str += "set location='" + string(location) + "' "
			flag = 1
		} else {
			str += ", location='" + string(location) + "' "
		}
	}
	if start != 0 {
		if flag == 0 {
			str += "set start_time=" + sstart + " "
			flag = 1
		} else {
			str += ", start_time=" + sstart + " "
		}
	}
	if end != 0 {
		if flag == 0 {
			str += "set end_time=" + send + " "
			flag = 1
		} else {
			str += ", end_time=" + send + " "
		}
	}
	str += "where id=" + sid
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
	sid := fmt.Sprintf("%d", id)
	if name != "" {
		str += "set name='" + name + "' "
	}
	str += "where id=" + sid
	_, err := globalVariable.DB.Exec(str)
	return err
}
