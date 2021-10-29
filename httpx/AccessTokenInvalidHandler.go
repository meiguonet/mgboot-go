package httpx

import (
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
)

type accessTokenInvalidHandler struct {
}

func NewAccessTokenInvalidHandler() *accessTokenInvalidHandler {
	return &accessTokenInvalidHandler{}
}

func (h *accessTokenInvalidHandler) GetExceptionName() string {
	return "builtin.AccessTokenInvalidException"
}

func (h *accessTokenInvalidHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.AccessTokenInvalidException); ok {
		return true
	}

	return false
}

func (h *accessTokenInvalidHandler) HandleException(err error) ResponsePayload {
	ex, ok := err.(BuiltinException.AccessTokenInvalidException)

	if !ok {
		return NewHttpError(500)
	}

	map1 := map[string]interface{}{
		"code": 1002,
		"msg":  ex.Error(),
		"data": nil,
	}

	return NewJsonResponse(map1)
}
