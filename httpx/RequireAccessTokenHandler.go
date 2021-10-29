package httpx

import (
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
)

type requireAccessTokenHandler struct {
}

func NewRequireAccessTokenHandler() *requireAccessTokenHandler {
	return &requireAccessTokenHandler{}
}

func (h *requireAccessTokenHandler) GetExceptionName() string {
	return "builtin.RequireAccessTokenException"
}

func (h *requireAccessTokenHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.RequireAccessTokenException); ok {
		return true
	}

	return false
}

func (h *requireAccessTokenHandler) HandleException(err error) ResponsePayload {
	ex, ok := err.(BuiltinException.RequireAccessTokenException)

	if !ok {
		return NewHttpError(500)
	}

	map1 := map[string]interface{}{
		"code": 1001,
		"msg":  ex.Error(),
		"data": nil,
	}

	return NewJsonResponse(map1)
}
