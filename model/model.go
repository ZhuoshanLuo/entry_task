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

type JoinOrExitRequest struct {
	SessionId  string `json:"sessionId"`
	ActivityId string `json:"activityId"`
	Action     string `json:"action"`
}

type Form struct {
	ActId    string `json:"activityId"`
	UserId   uint   `json:"userId"`
	JoinedAt uint   `json:"joinedAt"`
}

type ShowJoinedActivitiesRequest struct {
	SessionId string `json:"sessionId"`
}

type ShowJoinedActivitiesResponse struct {
	Title string `json:"title"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
}

type ActivityUserListRequest struct {
	ActivityId string `json:"activityId"`
}

type ActivityUserListResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type ActivityInfoRequest struct {
	SessionId uint `json:"sessionId"`
	ActivityId uint `json:"sessionId"`
}

type ActivityDetail struct{
	Title string `json:"title"`
	Start uint `json:"start"`
	End uint `json:"end"`
	Location string `json:"location"`
	Content string `json:"content"`
}

type ActivityInfoResponse struct {
	*ActivityDetail
	JoinStatus bool `json:"joinStatus"`
	UserList []ActivityUserListResponse
	CommentList []CommentListResponse
}