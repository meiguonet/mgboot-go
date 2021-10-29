package BuiltinMiddleware

import (
	"fmt"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	"github.com/meiguonet/mgboot-go/httpx"
	"github.com/meiguonet/mgboot-go/logx"
)

type requestLogMiddleware struct {
}

func NewRequestLogMiddleware() *requestLogMiddleware {
	return &requestLogMiddleware{}
}

func (m *requestLogMiddleware) GetName() string {
	return "builtin.RequestLog"
}

func (m *requestLogMiddleware) GetType() int {
	return MiddlewareType.PreHandle
}

func (m *requestLogMiddleware) GetOrder() int {
	return MiddlewareOrder.Highest
}

func (m *requestLogMiddleware) PreHandle(req *httpx.Request, resp *httpx.Response) {
	if !req.Next() || resp.HasError() {
		return
	}

	httpMethod := req.GetMethod()
	requestUrl := req.GetRequestUrl(true)
	clientIp := req.GetClientIp()

	if httpMethod == "" || requestUrl == "" || clientIp == "" {
		return
	}

	logger := logx.GetRequestLogLogger()
	msg := fmt.Sprintf("%s %s from %s", httpMethod, requestUrl, clientIp)
	logger.Info(msg)
	buf := req.GetRawBody()

	if len(buf) > 0 {
		logger.Debug(string(buf))
	}
}

func (m *requestLogMiddleware) PostHandle(_ *httpx.Request, _ *httpx.Response) {
}
