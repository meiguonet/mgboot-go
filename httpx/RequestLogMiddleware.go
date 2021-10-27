package httpx

import (
	"fmt"
	"github.com/meiguonet/mgboot-go"
	"github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
)

type RequestLogMiddleware struct {
	enabled bool
	logger  logx.Logger
}

func NewRequestLogMiddleware() RequestLogMiddleware {
	logger := mgboot.RequestLogLogger()
	var enabled bool

	if logger != nil {
		enabled = true
	}

	return RequestLogMiddleware{
		enabled: enabled,
		logger:  logger,
	}
}

func (m RequestLogMiddleware) GetType() int {
	return MiddlewareType.PreHandle
}

func (m RequestLogMiddleware) GetOrder() int {
	return MiddlewareOrder.Highest
}

func (m RequestLogMiddleware) PreHandle(req Request, resp Response) {
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

func (m RequestLogMiddleware) PostHandle(_ Request, _ Response) {
}
