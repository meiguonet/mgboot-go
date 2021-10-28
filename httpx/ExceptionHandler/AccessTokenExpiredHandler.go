package BuiltinExceptionHandler

import (
	"github.com/meiguonet/mgboot-go"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltintResponse "github.com/meiguonet/mgboot-go/httpx/response"
)

type accessTokenExpiredHandler struct {
}

func NewAccessTokenExpiredHandler() *accessTokenExpiredHandler {
	return &accessTokenExpiredHandler{}
}

func (h *accessTokenExpiredHandler) GetExceptionName() string {
	return mgboot.GetBuiltintExceptionName("AccessTokenExpired")
}

func (h *accessTokenExpiredHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.AccessTokenExpiredException); ok {
		return true
	}

	return false
}

func (h *accessTokenExpiredHandler) HandleException(err error) httpx.ResponsePayload {
	ex, ok := err.(BuiltinException.AccessTokenExpiredException)

	if !ok {
		return BuiltintResponse.NewHttpError(500)
	}

	map1 := map[string]interface{}{
		"code": 1003,
		"msg":  ex.Error(),
	}

	return BuiltintResponse.NewJsonResponse(map1)
}
