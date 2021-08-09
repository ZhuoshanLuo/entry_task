package svc

import (
	"github.com/ZhuoshanLuo/entry_task/codes"
	"github.com/ZhuoshanLuo/entry_task/model"
	"github.com/gin-gonic/gin"
)

/*
参数处理和权限处理文件
*/

func register(c *gin.Context) model.Response {
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

func login(c *gin.Context) model.Response {
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

func showJoinedActivities(c *gin.Context) model.Response {
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

func manageRegister(c *gin.Context) model.Response {
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
	//没有其他必要传入参数

	return doShowActivityType(userId)
}

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
	//业务逻辑
	return doShowAllUsers(userId)
}
