package model

//用户api请求参数
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

//显示所有活动函数请求参数
type ShowActivityRequest struct {
	SessionId uint `json:"sessionId"`
	Page      uint `json:"page"`
}

//用户加入的所有活动请求参数
type ShowJoinedActivitiesRequest struct {
	SessionId uint `json:"sessionId"`
}

//用户加入或退出活动请求参数
type JoinOrExitRequest struct {
	SessionId  uint `json:"sessionId"`
	ActivityId uint `json:"activityId"`
	Action     uint `json:"action"`
}

//显示活动详情请求参数
type ActivityInfoRequest struct {
	SessionId   uint `json:"sessionId"`
	ActivityId  uint `json:"activityId"`
	CommentPage uint `json:"page"`
}

//活动过滤器请求参数
type ActivitySelectorRequest struct {
	SessionId uint   `json:"sessionId"`
	Type      string `json:"type"`
	Start     uint   `json:"start"`
	End       uint   `json:"end"`
	Page      uint   `json:"page"`
}

//发表评论函数请求参数
type CreateCommentRequest struct {
	SessionId  uint   `json:"sessionId"`
	ActivityId uint   `json:"activityId"`
	Content    string `json:"content"`
}

//运营后台api请求参数
//添加活动函数的请求参数
type AddActivityRequest struct {
	SessionId uint   `json:"sessionId"`
	TypeId    uint   `json:"typeId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Location  string `json:"location"`
	Start     uint   `json:"start"`
	End       uint   `json:"End"`
}

//删除活动函数的请求参数
type DelActivityRequest struct {
	SessionId uint `json:"sessionId"`
	ActId     uint `json:"activityId"`
}

//编辑活动函数的请求参数
type EditActivityRequest struct {
	SessionId uint   `json:"sessionId"`
	Id        uint   `json:"id"`
	TypeId    uint   `json:"typeId"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Location  string `json:"location"`
	Start     uint   `json:"start"`
	End       uint   `json:"End"`
}

//添加活动类型的请求参数
type AddActivityTypeRequest struct {
	SessionId uint   `json:"sessionId"`
	TypeName  string `json:"typeName"`
}

//删除活动类型的请求参数
type DelActivityTypeRequest struct {
	SessionId uint `json:"sessionId"`
	Id        uint `json:"id"`
}

type ShowActivityTypeRequest struct {
	SessoinId uint `json:"sessionId"`
}

type ShowAllUsersRequest struct {
	SessionId uint `json:"sessionId"`
}

type EditActivityTypeRequest struct {
	SessionId uint   `json:"sessionId"`
	Id        uint   `json:"id"`
	Name      string `json:"name"`
}
