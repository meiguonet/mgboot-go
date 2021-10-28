package BuiltinMiddleware

import (
	"fmt"
	"github.com/meiguonet/mgboot-go"
	"github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	"github.com/meiguonet/mgboot-go/httpx"
)

type requestLogMiddleware struct {
	enabled bool
	logger  logx.Logger
}

func NewRequestLogMiddleware() *requestLogMiddleware {
	logger := mgboot.RequestLogLogger()
	var enabled bool

	if logger != nil {
		enabled = true
	}

	return &requestLogMiddleware{
		enabled: enabled,
		logger:  logger,
	}
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
	if !m.enabled || !req.Next() || resp.HasError() {
		return
	}

	httpMethod := req.GetMethod()
	requestUrl := req.GetRequestUrl(true)
	clientIp := req.GetClientIp()

	if httpMethod == "" || requestUrl == "" || clientIp == "" {
		return
	}

	msg := fmt.Sprintf("%s %s from %s", httpMethod, requestUrl, clientIp)
	m.logger.Info(msg)
	buf := req.GetRawBody()

	if len(buf) > 0 {
		m.logger.Debug(string(buf))
	}
}

func (m *requestLogMiddleware) PostHandle(_ *httpx.Request, _ *httpx.Response) {
}
