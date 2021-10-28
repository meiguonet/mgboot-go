package BuiltinExceptionHandler

import (
	"github.com/meiguonet/mgboot-go"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltintResponse "github.com/meiguonet/mgboot-go/httpx/response"
)

type requireAccessTokenHandler struct {
}

func NewRequireAccessTokenHandler() *requireAccessTokenHandler {
	return &requireAccessTokenHandler{}
}

func (h *requireAccessTokenHandler) GetExceptionName() string {
	return mgboot.GetBuiltintExceptionName("RequireAccessToken")
}

func (h *requireAccessTokenHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.RequireAccessTokenException); ok {
		return true
	}

	return false
}

func (h *requireAccessTokenHandler) HandleException(err error) httpx.ResponsePayload {
	ex, ok := err.(BuiltinException.RequireAccessTokenException)

	if !ok {
		return BuiltintResponse.NewHttpError(500)
	}

	map1 := map[string]interface{}{
		"code": 1001,
		"msg":  ex.Error(),
	}

	return BuiltintResponse.NewJsonResponse(map1)
}
