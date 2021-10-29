package httpx

import (
	"github.com/meiguonet/mgboot-go-common/util/jsonx"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
)

type validateExceptionHandler struct {
}

func NewValidateExceptionHandler() *validateExceptionHandler {
	return &validateExceptionHandler{}
}

func (h *validateExceptionHandler) GetExceptionName() string {
	return "builtin.ValidateException"
}

func (h *validateExceptionHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.ValidateException); ok {
		return true
	}

	return false
}

func (h *validateExceptionHandler) HandleException(err error) ResponsePayload {
	ex, ok := err.(BuiltinException.ValidateException)

	if !ok {
		return NewHttpError(500)
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
		"data": nil,
	}

	return NewJsonResponse(map1)
}
