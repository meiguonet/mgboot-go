package BuiltinExceptionHandler

import (
	"github.com/meiguonet/mgboot-go"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltintResponse "github.com/meiguonet/mgboot-go/httpx/response"
)

type accessTokenInvalidHandler struct {
}

func NewAccessTokenInvalidHandler() *accessTokenInvalidHandler {
	return &accessTokenInvalidHandler{}
}

func (h *accessTokenInvalidHandler) GetExceptionName() string {
	return mgboot.GetBuiltintExceptionName("AccessTokenInvalid")
}

func (h *accessTokenInvalidHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.AccessTokenInvalidException); ok {
		return true
	}

	return false
}

func (h *accessTokenInvalidHandler) HandleException(err error) httpx.ResponsePayload {
	ex, ok := err.(BuiltinException.AccessTokenInvalidException)

	if !ok {
		return BuiltintResponse.NewHttpError(500)
	}

	map1 := map[string]interface{}{
		"code": 1002,
		"msg":  ex.Error(),
	}

	return BuiltintResponse.NewJsonResponse(map1)
}
