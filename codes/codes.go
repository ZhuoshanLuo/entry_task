package codes

import "net/http"

type Code uint32

const (
	OK   Code = 0
	Fail Code = 1
)

func HTTPStatusFromCode(code Code) int {
	codes := map[Code]int{
		OK:   http.StatusOK,
		Fail: http.StatusBadRequest,
	}
	return codes[code]
}

func Errorf(code Code) string {
	var codeMsg string
	switch code {
	case OK:
		codeMsg = "SUCCESS!"
	case Fail:
		codeMsg = "FAIL!"
	}
	return codeMsg
}
