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
