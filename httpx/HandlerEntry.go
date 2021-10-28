package httpx

import (
	"github.com/meiguonet/mgboot-go/mvc"
)

type ActionFunc func(req *Request, resp *Response) (ResponsePayload, error)

type HandlerEntry struct {
	routeRule   *mvc.RouteRule
	actionFunc  ActionFunc
}

func NewHandlerEntry(routeRule *mvc.RouteRule) *HandlerEntry {
	entry := &HandlerEntry{routeRule: routeRule}
	return entry
}

func (e *HandlerEntry) GetRouteRule() *mvc.RouteRule {
	return e.routeRule
}

func (e *HandlerEntry) WithActionFunc(fn ActionFunc) *HandlerEntry {
	e.actionFunc = fn
	return e
}

func (e *HandlerEntry) HandleRequest(req *Request, resp *Response) {
	var payload ResponsePayload
	var err error

	defer func() {
		if ex, ok := recover().(error); ok {
			err = ex
		}
	}()

	payload, err = e.actionFunc(req, resp)

	if err != nil {
		resp.WithError(err)
	} else {
		resp.WithPayload(payload)
	}

	resp.Send()
}
