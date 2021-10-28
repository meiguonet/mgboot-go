package BuiltinExceptionHandler

import (
	"github.com/meiguonet/mgboot-go-common/util/jsonx"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltintResponse "github.com/meiguonet/mgboot-go/httpx/response"
)

type validateHandler struct {
}

func NewValidateHandler() *validateHandler {
	return &validateHandler{}
}

func (h *validateHandler) GetExceptionName() string {
	return "builtin.ValidateException"
}

func (h *validateHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.ValidateException); ok {
		return true
	}

	return false
}

func (h *validateHandler) HandleException(err error) httpx.ResponsePayload {
	ex, ok := err.(BuiltinException.ValidateException)

	if !ok {
		return BuiltintResponse.NewHttpError(500)
	}

	var msg string

	if ex.IsFailfast() {
		msg = ex.Error()
	} else {
		msg = jsonx.ToJson(ex.GetValidateErrors())
	}

	map1 := map[string]interface{}{
		"code": 1006,
		"msg":  msg,
	}

	return BuiltintResponse.NewJsonResponse(map1)
}
