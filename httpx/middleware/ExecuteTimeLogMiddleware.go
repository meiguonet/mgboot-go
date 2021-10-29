package BuiltinMiddleware

import (
	"fmt"
	"github.com/meiguonet/mgboot-go-common/util/numberx"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	"github.com/meiguonet/mgboot-go/httpx"
	"github.com/meiguonet/mgboot-go/logx"
	"time"
)

type executeTimeLogMiddleware struct {
}

func NewExecuteTimeLogMiddleware() *executeTimeLogMiddleware {
	return &executeTimeLogMiddleware{}
}

func (m *executeTimeLogMiddleware) GetName() string {
	return "builtin.ExecuteTimeLog"
}

func (m *executeTimeLogMiddleware) GetType() int {
	return MiddlewareType.PostHandle
}

func (m *executeTimeLogMiddleware) GetOrder() int {
	return MiddlewareOrder.Lowest
}

func (m *executeTimeLogMiddleware) PreHandle(_ *httpx.Request, _ *httpx.Response) {
}

func (m *executeTimeLogMiddleware) PostHandle(req *httpx.Request, resp *httpx.Response) {
	httpMethod := req.GetMethod()
	requestUrl := req.GetRequestUrl(true)

	if httpMethod == "" || requestUrl == "" {
		return
	}

	elapsedTime := m.calcElapsedTime(req)
	msg := fmt.Sprintf("%s %s, , total elapsed time: %s", httpMethod, requestUrl, elapsedTime)
	logx.GetExecuteTimeLogLogger().Info(msg)
	resp.WithExtraHeader("X-Response-Time", elapsedTime)
}

func (m *executeTimeLogMiddleware) calcElapsedTime(req *httpx.Request) string {
	d1 := time.Now().Sub(req.GetExecStart())

	if d1 < time.Second {
		return fmt.Sprintf("%dms", d1)
	}

	n1 := d1.Seconds()
	return numberx.ToDecimalString(n1, 3) + "ms"
}
