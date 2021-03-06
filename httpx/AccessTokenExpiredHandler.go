package httpx

import BuiltinException "github.com/meiguonet/mgboot-go/exception"

type accessTokenExpiredHandler struct {
}

func NewAccessTokenExpiredHandler() *accessTokenExpiredHandler {
	return &accessTokenExpiredHandler{}
}

func (h *accessTokenExpiredHandler) GetExceptionName() string {
	return "builtin.AccessTokenExpiredException"
}

func (h *accessTokenExpiredHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.AccessTokenExpiredException); ok {
		return true
	}

	return false
}

func (h *accessTokenExpiredHandler) HandleException(err error) ResponsePayload {
	ex, ok := err.(BuiltinException.AccessTokenExpiredException)

	if !ok {
		return NewHttpError(500)
	}

	map1 := map[string]interface{}{
		"code": 1003,
		"msg":  ex.Error(),
		"data": nil,
	}

	return NewJsonResponse(map1)
}
