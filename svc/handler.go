package svc

import (
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/gin-gonic/gin"
)

/*
参数处理和权限处理文件
*/
//注册
func register(c *gin.Context) model.Response {
	//参数校验
	isValid, req := registerParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doRegister(*req)
}

//登陆
func login(c *gin.Context) model.Response {
	//参数校验
	isValid, req := loginParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doLogin(*req)
}

//显示活动
func showActivities(c *gin.Context) model.Response {
	//用户是否登陆，如果登陆，从cookie中得到user id
	userId, err := CheckUserLogin(c)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	//参数校验
	isValid, req := showActivitiesParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doShowActivities(*req, userId)
}

//活动筛选器
func activitiesSelector(c *gin.Context) model.Response {
	//用户是否登陆，如果登陆，从cookie中得到user id
	userId, err := CheckUserLogin(c)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	//参数校验
	isValid, req := activitiesSelectorParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doActivitiesSelector(*req, userId)
}

//创建评论
func createComment(c *gin.Context) model.Response {
	//用户是否登陆，如果登陆，从cookie中得到user id
	userId, err := CheckUserLogin(c)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	if userId == 0 {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := createCommentParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doCreateComment(*req, userId)
}

//显示加入的所有活动
func showJoinedActivities(c *gin.Context) model.Response {
	//用户是否登陆，如果登陆，从cookie中得到user id，本接口一定要在登陆态
	userId, err := CheckUserLogin(c)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	if userId == 0 {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//没有其他必要传入的参数，不用做参数校验

	return doShowJoinedActivities(userId)
}

//加入或退出活动
func joinOrExit(c *gin.Context) model.Response {
	//用户是否登陆，如果登陆，从cookie中得到user id
	userId, err := CheckUserLogin(c)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}
	if userId == 0 {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := joinOrExitParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doJoinOrExit(*req, userId)
}

//显示活动信息
func activityInfo(c *gin.Context) model.Response {
	//用户是否登陆，如果登陆，从cookie中得到user id
	userId, err := CheckUserLogin(c)
	if err != nil {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.MysqlError,
				Msg:  codes.Errorf(codes.MysqlError),
			},
		}
	}

	//参数校验
	isValid, req := activityInfoParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doActivityInfo(*req, userId)
}

//注册
func manageRegister(c *gin.Context) model.Response {
	//参数校验
	isValid, req := manageRegisterParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doManageRegister(*req)
}

//登陆
func manageLogin(c *gin.Context) model.Response {
	//参数校验
	isValid, req := manageLoginParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	return doManageLogin(*req)
}

//添加活动
func addActivity(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	//参数校验
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	isValid, req := addActivityParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doAddActivity(*req, userId)
}

//删除活动
func delActivity(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	//参数校验
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	isValid, req := delActivityParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doDelActivity(*req, userId)
}

//编辑活动
func editActivity(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := editActivityParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doEditActivity(*req, userId)
}

//添加活动类型
func addActivityType(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := addActivityTypeParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doAddActivityType(*req, userId)
}

//显示活动类型
func showActivityType(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}

	isValid, req := ShowActivityTypeRequestParameterValid(c)

	return doShowActivityType(*req, userId)
}

//编辑活动类型
func editActivityType(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := editActivityTypeParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doEditActivityType(*req, userId)
}

//删除活动类型
func delActivityType(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := delActivityTypeParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}

	return doDelActivityType(*req, userId)
}

//显示所有用户
func showAllUsers(c *gin.Context) model.Response {
	//身份验证
	userId, isValid := CheckIdentity(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.Forbidden,
				Msg:  codes.Errorf(codes.Forbidden),
			},
		}
	}
	//参数校验
	isValid, req := showAllUsersParameterValid(c)
	if isValid == false {
		return model.Response{
			Status: model.CodeMsg{
				Code: codes.ParameterError,
				Msg:  codes.Errorf(codes.ParameterError),
			},
		}
	}
	//业务逻辑
	return doShowAllUsers(*req, userId)
}
