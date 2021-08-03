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

type Session struct {
	Id     uint `json:"id"`
	UserId uint `json:"userId"`
}

//结构体
type Response struct {
	Code codes.Code  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

//注册函数请求参数
type RegisterRequest struct {
	Name   string `json:"name"`
	Passwd string `json:"Passwd"`
	Email  string `json:"email"`
	Avatar string `json:"image"`
}

//登陆函数请求参数
type LoginRequest struct {
	Name   string `json:"name"`
	Passwd string `json:"Passwd"`
}

//显示所有活动函数请求参数、返回参数
type ShowActivityRequest struct {
	SessionId uint `json:"sessionId"`
	Page      uint   `json:"page"`
}

type ShowActivitiesResponse struct {
	Title      string `json:"title"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	JoinStatus bool   `json:"status"`
}

//活动过滤器请求参数、返回参数
type ActivitySelectorRequest struct {
	SessionId uint `json:"sessionId"`
	Type      string `json:"type"`
	Start     uint `json:"start"`
	End       uint `json:"end"`
	Page      uint   `json:"page"`
}

type ActivitySelectorResponse struct {
	Title      string `json:"title"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	JoinStatus bool   `json:"joinStatus"`
}

//发表评论函数请求参数
type CreateCommentRequest struct {
	SessionId  uint `json:"sessionId"`
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
}

//用户加入的所有活动请求参数、返回参数
type ShowJoinedActivitiesRequest struct {
	SessionId uint `json:"sessionId"`
}

type ShowJoinedActivitiesResponse struct {
	Title string `json:"title"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
}

//用户加入或退出活动请求参数
type JoinOrExitRequest struct {
	SessionId  uint `json:"sessionId"`
	ActivityId uint `json:"activityId"`
	Action     uint `json:"action"`
}

//显示活动详情请求参数、返回参数
type ActivityInfoRequest struct {
	SessionId uint `json:"sessionId"`
	ActivityId uint `json:"sessionId"`
	CommentPage uint `json:"commentPage"`
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



type CommentListResponse struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt uint   `json:"createdAt"`
}

type ActivityUserListResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}




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
	UserId   uint   `json:"userId"`
	JoinedAt uint   `json:"joinedAt"`
}