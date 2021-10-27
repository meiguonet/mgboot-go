package httpx

import (
	"fmt"
	"github.com/meiguonet/mgboot-go"
	"github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go-common/util/numberx"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	"time"
)

type ExecuteTimeLogMiddleware struct {
	enabled bool
	logger  logx.Logger
}

func NewExecuteTimeLogMiddleware() ExecuteTimeLogMiddleware {
	logger := mgboot.ExecuteTimeLogLogger()
	var enabled bool

	if logger != nil {
		enabled = true
	}

	return ExecuteTimeLogMiddleware{
		enabled: enabled,
		logger:  logger,
	}
}

func (m ExecuteTimeLogMiddleware) GetType() int {
	return MiddlewareType.PostHandle
}

func (m ExecuteTimeLogMiddleware) GetOrder() int {
	return MiddlewareOrder.Lowest
}

func (m ExecuteTimeLogMiddleware) PreHandle(_ Request, _ Response) {
}

func (m ExecuteTimeLogMiddleware) PostHandle(req Request, resp Response) {
	if !m.enabled || resp.HasError() {
		return
	}

	httpMethod := req.GetMethod()
	requestUrl := req.GetRequestUrl(true)

	if httpMethod == "" || requestUrl == "" {
		return
	}

	elapsedTime := m.calcElapsedTime(req)
	msg := fmt.Sprintf("%s %s, , total elapsed time: %s", httpMethod, requestUrl, elapsedTime)
	m.logger.Info(msg)
	resp.WithExtraHeader("X-Response-Time", elapsedTime)
}

func (m ExecuteTimeLogMiddleware) calcElapsedTime(req Request) string {
	d1 := time.Now().Sub(req.GetExecStart())

	if d1 < time.Second {
		return fmt.Sprintf("%dms", d1)
	}

	n1 := d1.Seconds()
	return numberx.ToDecimalString(n1, 3) + "ms"
}
