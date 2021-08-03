package codes

import "net/http"

type Code uint32

const (
	OK            Code = 0
	Fail          Code = 1
	MissParameter Code = 2
	UserExist     Code = 3
	MysqlError    Code = 4
	UserNotExist  Code = 5
	PassWordError Code = 6
	NotLogin      Code = 7
	BindJsonError Code = 8
	Forbidden     Code = 9
)

func HTTPStatusFromCode(code Code) int {
	codes := map[Code]int{
		OK:            http.StatusOK,
		Fail:          http.StatusBadRequest,
		MissParameter: http.StatusBadRequest,
		UserExist:     http.StatusBadRequest,
		MysqlError:    http.StatusInternalServerError,
		UserNotExist:  http.StatusBadRequest,
		PassWordError: http.StatusBadRequest,
		NotLogin:      http.StatusBadRequest,
		BindJsonError: http.StatusBadRequest,
		Forbidden:     http.StatusForbidden,
	}
	return codes[code]
}

func Errorf(code Code) string {
	var codeMsg string
	switch code {
	case OK:
		codeMsg = "Success!"
	case Fail:
		codeMsg = "Fail!"
	case MissParameter:
		codeMsg = "Parameters missing or format error!"
	case UserExist:
		codeMsg = "User have exist!"
	case MysqlError:
		codeMsg = "Accessing the database error!"
	case UserNotExist:
		codeMsg = "User is not exist!"
	case PassWordError:
		codeMsg = "Password error!"
	case NotLogin:
		codeMsg = "User not login!"
	case BindJsonError:
		codeMsg = "Bind json error!"
	case Forbidden:
		codeMsg = "You have no authority！"
	}
	return codeMsg
}
