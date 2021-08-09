package model

//通用api
type RegisterMsg struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

//用户api请求参数
//登陆函数请求参数
type LoginMsg struct {
	Name   string `json:"name"`
	Passwd string `json:"passwd"`
}

//显示所有活动函数请求参数
type ShowActivtyRequest struct {
	Page uint `json:"page"`
}

//用户加入或退出活动请求参数
type JoinOrExitRequest struct {
	ActivityId uint `json:"activity_id"`
	Action     uint `json:"action"`
}

//显示活动详情请求参数
type ActivityInfoRequest struct {
	ActivityId  uint `json:"activity_id"`
	CommentPage uint `json:"comment_page"`
}

//活动过滤器请求参数
type ActivitySelectorRequest struct {
	Type  string `json:"type"`
	Start uint   `json:"start"`
	End   uint   `json:"end"`
	Page  uint   `json:"page"`
}

//发表评论函数请求参数
type CreateCommentRequest struct {
	ActivityId uint   `json:"activity_id"`
	Content    string `json:"content"`
}

//运营后台api请求参数
//添加活动函数的请求参数
type AddActivityRequest struct {
	TypeId   uint   `json:"type_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Location string `json:"location"`
	Start    uint   `json:"start"`
	End      uint   `json:"end"`
}

//删除活动函数的请求参数
type DelActivityRequest struct {
	ActId uint `json:"act_id"`
}

//编辑活动函数的请求参数
type EditActivityRequest struct {
	Id       uint   `json:"id"`
	TypeId   uint   `json:"typeId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	Location string `json:"location"`
	Start    uint   `json:"start"`
	End      uint   `json:"End"`
}

//添加活动类型的请求参数
type AddActivityTypeRequest struct {
	TypeName string `json:"type_name"`
}

//删除活动类型的请求参数
type DelActivityTypeRequest struct {
	Id uint `json:"id"`
}

type EditActivityTypeRequest struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

type ShowActivityTypeRequest struct {
	Page uint `json:"page"`
}

type ShowAllUsersRequest struct {
	Page uint `json:"page"`
}
