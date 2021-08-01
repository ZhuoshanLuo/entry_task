package model

import "example.com/greetings/codes"

//结构体
type Response struct {
	Code codes.Code  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type Session struct {
	SessionId uint `json:"sessionId"`
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
	Id        uint   `json:"id"`
	Passwd    string `json:"passwd"`
	IsAdmin   bool   `json:"IsAdmin"`
	CreatedAt uint   `json:"createdAt"`
	UserMsg
}

type CommentMsg struct {
	UserName string `json:"username"`
	Time     uint   `json:"time"`
	Content  string `json:"content"`
}

type Comment struct {
	Id         uint   `json:"id"`
	UserId     uint   `json:"userId"`
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
	CreatedAt  uint   `json:"CreatedAt"`
}

//显示所有活动的返回值
type ShowActivitiesRes struct {
	Activity
	JoinStatus uint `json:"status"`
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
	Passwd string `json:"passwd"`
	Email  string `json:"email"`
	Avatar string `json:"image"`
}

type LoginRequest struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

type ShowActivityRequest struct {
	UserId string `json:"userId"`
}

type CreateCommentRequest struct {
	Session
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
}
