package model

import "example.com/greetings/codes"

type ShowActivitiesResponse struct {
	Title      string `json:"title"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	JoinStatus bool   `json:"status"`
}

type ActivitySelectorResponse struct {
	Title      string `json:"title"`
	Start      uint   `json:"start"`
	End        uint   `json:"end"`
	JoinStatus bool   `json:"joinStatus"`
}

type ShowJoinedActivitiesResponse struct {
	Title string `json:"title"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
}

type ActivityInfoResponse struct {
	*ActivityDetail
	JoinStatus  bool `json:"joinStatus"`
	UserList    []ActivityUserListResponse
	CommentList []CommentListResponse
}

type CommentListResponse struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	CreatedAt uint   `json:"createdAt"`
}

type ActivityUserListResponse struct {
	Name   string `db:"name"`
	Avatar string `db:"avatar"`
}

//结构体
type Response struct {
	Code codes.Code  `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type ActivityDetail struct {
	Title    string `db:"title"`
	Start    uint   `db:"start_time"`
	End      uint   `db:"end_time"`
	Location string `db:"location"`
	Content  string `db:"content"`
}

type ShowActivityTypeResponse struct {
	TypeName string `json:"typeName"`
}

type ShowAllUsersResponse struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Avatart string `json:"avatar"`
}
