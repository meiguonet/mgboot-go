package httpx

import (
	"github.com/meiguonet/mgboot-go/mvc"
	"net/http"
)

type ActionFunc func(req *Request, resp *Response) (ResponsePayload, error)
type HandlerFunc func(out http.ResponseWriter, req *http.Request, pathVariables map[string]string)

type HandlerEntry struct {
	routeRule   *mvc.RouteRule
	actionFunc  ActionFunc
	handlerFunc HandlerFunc
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

func (e *HandlerEntry) GetActionFunc() ActionFunc {
	return e.actionFunc
}

func (e *HandlerEntry) WithHandlerFunc(fn HandlerFunc) *HandlerEntry {
	e.handlerFunc = fn
	return e
}

func (e *HandlerEntry) HandleRequest(out http.ResponseWriter, req *http.Request, pathVariables map[string]string) {
	e.handlerFunc(out, req, pathVariables)
}
