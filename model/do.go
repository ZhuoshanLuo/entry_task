package model

//表数据结构
type User struct {
	Name      string `json:"name"`
	Passwd    string `json:"Passwd"`
	Email     string `json:"email"`
	Avatar    string `json:"image"`
	IsAdmin   bool   `json:"IsAdmin"`
	CreatedAt uint   `json:"createdAt"`
}

type Comment struct {
	UserId     uint   `json:"userId"`
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
	CreatedAt  uint   `json:"CreatedAt"`
}

type Form struct {
	ActId    uint `json:"activityId"`
	UserId   uint `json:"userId"`
	JoinedAt uint `json:"joinedAt"`
}

type Session struct {
	Id     uint `json:"id"`
	UserId uint `json:"userId"`
}

type Activity struct {
	TypeId   uint   `json:"typeId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Location string `json:"location"`
	Start    uint   `json:"start"`
	End      uint   `json:"end"`
}

type ActivityType struct {
	Name string `json:"name"`
}