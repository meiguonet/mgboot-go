package httpx

import (
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
)

type unkownErrorHandler struct {
}

func NewUnkownErrorHandler() *unkownErrorHandler {
	return &unkownErrorHandler{}
}

func (h *unkownErrorHandler) GetExceptionName() string {
	return "built.UnkownErrorException"
}

func (h *unkownErrorHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.UnkownErrorException); ok {
		return true
	}

	return false
}

func (h *unkownErrorHandler) HandleException(err error) ResponsePayload {
	var msg string
	ex, ok := err.(BuiltinException.UnkownErrorException)

	if ok {
		msg = ex.Error()
	} else {
		msg = "unknow error found"
	}

	map1 := map[string]interface{}{
		"code": 500,
		"msg":  msg,
		"data": nil,
	}

	return NewJsonResponse(map1)
}
