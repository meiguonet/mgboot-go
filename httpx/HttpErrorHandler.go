package httpx

import (
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
)

type httpErrorHandler struct {
}

func NewHttpErrorHandler() *httpErrorHandler {
	return &httpErrorHandler{}
}

func (h *httpErrorHandler) GetExceptionName() string {
	return "builtin.HttpErrorException"
}

func (h *httpErrorHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.HttpErrorException); ok {
		return true
	}

	return false
}

func (h *httpErrorHandler) HandleException(err error) ResponsePayload {
	ex, ok := err.(BuiltinException.HttpErrorException)

	if !ok {
		return NewHttpError(500)
	}

	return NewHttpError(ex.GetStatusCode())
}
