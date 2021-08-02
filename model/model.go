package model

import "example.com/greetings/codes"

type Config struct {
	Db struct {
		Driver   string `yaml:"driver"`
		SqlUser  string `yaml:"sqlUser"`
		Passwd   string `yaml:"passwd"`
		Host     string `yaml:"host"`
		Database string `yaml:"database"`
	}
}

//结构体
type Response struct {
	Code codes.Code  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Session struct {
	Id     uint `json:"id"`
	UserId uint `json:"userId"`
}

type Activity struct {
	Title string `json:"title"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
}

type ActivityDetail struct {
	Content  string   `json:"content"`
	Location string   `json:"location"`
	Act      Activity `json:"activity"`
}

type UserMsg struct {
	Name   string `json:"name"`
	Avatar string `json:"image"`
	Email  string `json:"email"`
}

type User struct {
	Name      string `json:"name"`
	Passwd    string `json:"Passwd"`
	Email     string `json:"email"`
	Avatar    string `json:"image"`
	IsAdmin   bool   `json:"IsAdmin"`
	CreatedAt uint   `json:"createdAt"`
}

type CommentMsg struct {
	UserName string `json:"username"`
	Time     uint   `json:"time"`
	Content  string `json:"content"`
}

type Comment struct {
	UserId     uint   `json:"userId"`
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
	CreatedAt  uint   `json:"CreatedAt"`
}

//显示所有活动的返回值
type ShowActivitiesRes struct {
	Title      string `json:"title"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	JoinStatus bool   `json:"status"`
}

//活动选择器函数的返回值
type SelectorStruct struct {
	Activity
	Response
	JoinStatus uint `json:"status"'`
}

//活动详情的返回值
type ADetailStruct struct {
}

type RegisterRequest struct {
	Name   string `json:"name"`
	Passwd string `json:"Passwd"`
	Email  string `json:"email"`
	Avatar string `json:"image"`
}

type LoginRequest struct {
	Name   string `json:"name"`
	Passwd string `json:"Passwd"`
}

type ShowActivityRequest struct {
	SessionId string `json:"sessionId"`
	Page      uint   `json:"page"`
}

type CreateCommentRequest struct {
	SessionId  string `json:"sessionId"`
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
}

type ActivitySelectorRequest struct {
	SessionId string `json:"sessionId"`
	Type      string `json:"type"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Page      uint   `json:"page"`
}

type ActivitySelectorResponse struct {
	Title      string `json:"title"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	JoinStatus bool   `json:"joinStatus"`
}

type CommentListRequest struct {
	ActivityId string `json:"activityId"`
	Page       uint   `json:"page"`
}

type CommentListResponse struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt uint   `json:"createdAt"`
}
