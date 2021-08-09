package model

import "github.com/ZhuoshanLuo/entry_task/codes"

//结构体

type CodeMsg struct {
	Code codes.Code `json:"code"`
	Msg  string     `json:"msg"`
}

type Response struct {
	Status     CodeMsg     `json:"status"`
	Data       interface{} `json:"data"`
	Session_id uint        `json:"session_id"`
}

type ActivityProfile struct {
	Title string `json:"title"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
}

type ShowJoinedActivityResponse struct {
	Activities []*ActivityProfile `json:"activities"`
}

type UserActivityInfo struct {
	ActivityProfile `json:"activity_profile"`
	JoinStatus      bool `json:"join_status"`
}

type ShowActivitiesResponse struct {
	Activities []*UserActivityInfo `json:"activities"`
}

type ActivitySelectorResponse struct {
	Activities []*UserActivityInfo `json:"activities"`
}

type ActivityInfoResponse struct {
	*ActivityDetail
	JoinStatus  bool `json:"join_status"`
	UserList    []UserPublicMsg
	CommentList []CommentListResponse
}

type CommentListResponse struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt uint   `json:"created_at"`
}

type ActivityDetail struct {
	Title    string `db:"title"`
	Start    uint   `db:"start_time"`
	End      uint   `db:"end_time"`
	Location string `db:"location"`
	Content  string `db:"content"`
}

type ShowActivityTypeResponse struct {
	TypeName string `json:"type_name"`
}

type UserPublicMsg struct {
	Name   string `db:"name"`
	Email  string `db:"email"`
	Avatar string `db:"avatar"`
}
