package model

//表数据结构
type User struct {
	Name      string `db:"name"`
	Passwd    string `db:"passwd"`
	Email     string `db:"email"`
	Avatar    string `db:"avatar"`
	IsAdmin   bool   `db:"is_admin"`
	CreatedAt uint   `db:"created_at"`
}

type Comment struct {
	UserId     uint   `db:"user_id"`
	ActivityId uint   `db:"activity_id"`
	Content    string `db:"content"`
	CreatedAt  uint   `db:"created_at"`
}

type Form struct {
	ActId    uint `db:"activity_id"`
	UserId   uint `db:"user_id"`
	JoinedAt uint `db:"joined_at"`
}

type Session struct {
	Id     uint `db:"id"`
	UserId uint `db:"user_id"`
}

type Activity struct {
	TypeId   uint   `db:"typeId"`
	Title    string `db:"title"`
	Content  string `db:"content"`
	Location string `db:"location"`
	Start    uint   `db:"start"`
	End      uint   `db:"end"`
}

type ActivityType struct {
	Name string `db:"name"`
}