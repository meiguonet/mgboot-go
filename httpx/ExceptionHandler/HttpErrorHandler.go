package BuiltinExceptionHandler

import (
	"github.com/meiguonet/mgboot-go"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltintResponse "github.com/meiguonet/mgboot-go/httpx/response"
)

type httpErrorHandler struct {
}

func NewHttpErrorHandler() *httpErrorHandler {
	return &httpErrorHandler{}
}

func (h *httpErrorHandler) GetExceptionName() string {
	return mgboot.GetBuiltintExceptionName("HttpError")
}

func (h *httpErrorHandler) MatchException(err error) bool {
	if _, ok := err.(BuiltinException.HttpErrorException); ok {
		return true
	}

	return false
}

func (h *httpErrorHandler) HandleException(err error) httpx.ResponsePayload {
	ex, ok := err.(BuiltinException.HttpErrorException)

	if !ok {
		return BuiltintResponse.NewHttpError(500)
	}

	return BuiltintResponse.NewHttpError(ex.GetStatusCode())
}
