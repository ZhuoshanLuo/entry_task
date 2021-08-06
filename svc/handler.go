package svc

import (
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/gin-gonic/gin"
)

func NewHandlerFunc(f Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		code, data := f(c)
		res := model.Response{
			Code: code,
			Msg:  codes.Errorf(code),
		}
		if data != nil {
			res.Data = data
		}
		c.JSON(codes.HTTPStatusFromCode(code), res)
	}
}

//用户接口
var (
	Register             = NewHandlerFunc(register)
	Login                = NewHandlerFunc(login)
	ShowActivities       = NewHandlerFunc(showActivities)
	ActivitiesSelector   = NewHandlerFunc(activitiesSelector)
	ActivityInfo         = NewHandlerFunc(activityInfo)
	CreateComment        = NewHandlerFunc(createComment)
	ShowJoinedActivities = NewHandlerFunc(showJoinedActivities)
	JoinOrExit           = NewHandlerFunc(joinOrExit)
)

//管理员接口
var (
	ManageRegister   = NewHandlerFunc(manageRegister)
	ManageLogin      = NewHandlerFunc(manageLogin)
	AddActivity      = NewHandlerFunc(addActivity)
	DelActivity      = NewHandlerFunc(delActivity)
	AddActivityType  = NewHandlerFunc(addActivityType)
	EditActivity     = NewHandlerFunc(editActivity)
	ShowAllUsers     = NewHandlerFunc(showAllUsers)
	EditActivityType = NewHandlerFunc(editActivityType)
	DelActivityType  = NewHandlerFunc(delActivityType)
	ShowActivityType = NewHandlerFunc(showActivityType)
)
