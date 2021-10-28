package mvc

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go/httpx"
)

type RouteRule struct {
	httpMethod        string
	requestMapping    string
	regex             string
	pathVariableNames []string
	jwtSettingsKey    string
	validateRules     []string
	failfast          bool
	handlerFuncArgs   []HandlerFuncArgInfo
	middlewares       []httpx.Middleware
}

func NewRouteRule(map1 map[string]interface{}) *RouteRule {
	pathVariableNames := make([]string, 0)

	if a1, ok := map1["pathVariableNames"].([]string); ok && len(a1) > 0 {
		pathVariableNames = a1
	}

	validateRules := make([]string, 0)

	if a1, ok := map1["validateRules"].([]string); ok && len(a1) > 0 {
		validateRules = a1
	}

	handlerFuncArgs := make([]HandlerFuncArgInfo, 0)

	if a1, ok := map1["handlerFuncArgs"].([]HandlerFuncArgInfo); ok && len(a1) > 0 {
		handlerFuncArgs = a1
	}

	middlewares := make([]httpx.Middleware, 0)

	if a1, ok := map1["middlewares"].([]httpx.Middleware); ok && len(a1) > 0 {
		middlewares = a1
	}

	return &RouteRule{
		httpMethod:        castx.ToString(map1["httpMethod"]),
		requestMapping:    castx.ToString(map1["requestMapping"]),
		regex:             castx.ToString(map1["regex"]),
		pathVariableNames: pathVariableNames,
		jwtSettingsKey:    castx.ToString(map1["jwtSettingsKey"]),
		validateRules:     validateRules,
		failfast:          castx.ToBool(map1["failfast"]),
		handlerFuncArgs:   handlerFuncArgs,
		middlewares:       middlewares,
	}
}

func (rr *RouteRule) HttpMethod() string {
	return rr.httpMethod
}

func (rr *RouteRule) RequestMapping() string {
	return rr.requestMapping
}

func (rr *RouteRule) Regex() string {
	return rr.regex
}

func (rr *RouteRule) PathVariableNames() []string {
	return rr.pathVariableNames
}

func (rr *RouteRule) JwtSettingsKey() string {
	return rr.jwtSettingsKey
}

func (rr *RouteRule) ValidateRules() []string {
	return rr.validateRules
}

func (rr *RouteRule) IsFailfast() bool {
	return rr.failfast
}

func (rr *RouteRule) HandlerFuncArgs() []HandlerFuncArgInfo {
	return rr.handlerFuncArgs
}

func (rr *RouteRule) WithMiddleware(m httpx.Middleware) *RouteRule {
	rr.middlewares = append(rr.middlewares, m)
	return rr
}

func (rr *RouteRule) WithMiddlewares(entries []httpx.Middleware) *RouteRule {
	if len(entries) > 0 {
		rr.middlewares = append(rr.middlewares, entries...)
	}

	return rr
}

func (rr *RouteRule) Middlewares() []httpx.Middleware {
	return rr.middlewares
}
