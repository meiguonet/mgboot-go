package BuiltinExceptionHandler

import (
	"github.com/meiguonet/mgboot-go-dal/dbx"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltintResponse "github.com/meiguonet/mgboot-go/httpx/response"
)

type dbExceptionHandler struct {
}

func NewDbExceptionHandler() *dbExceptionHandler {
	return &dbExceptionHandler{}
}

func (h *dbExceptionHandler) GetExceptionName() string {
	return "dbx.DbException"
}

func (h *dbExceptionHandler) MatchException(err error) bool {
	if _, ok := err.(dbx.DbException); ok {
		return true
	}

	return false
}

func (h *dbExceptionHandler) HandleException(_ error) httpx.ResponsePayload {
	return BuiltintResponse.NewHttpError(500)
}
