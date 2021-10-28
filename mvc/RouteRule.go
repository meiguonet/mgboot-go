package mvc

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
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

	return &RouteRule{
		httpMethod:        castx.ToString(map1["httpMethod"]),
		requestMapping:    castx.ToString(map1["requestMapping"]),
		regex:             castx.ToString(map1["regex"]),
		pathVariableNames: pathVariableNames,
		jwtSettingsKey:    castx.ToString(map1["jwtSettingsKey"]),
		validateRules:     validateRules,
		failfast:          castx.ToBool(map1["failfast"]),
		handlerFuncArgs:   handlerFuncArgs,
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
