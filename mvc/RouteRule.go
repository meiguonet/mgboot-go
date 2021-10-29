package mvc

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"time"
)

type RouteRule struct {
	httpMethod        string
	requestMapping    string
	regex             bool
	pathVariableNames []string
	jwtSettingsKey    string
	validateRules     []string
	failfast          bool
	rateLimitTotal    int
	rateLimitDuration time.Duration
	rateLimitByIp     bool
	customMiddlewares []string
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

	var rateLimitDuration time.Duration

	if d1, ok := map1["rateLimitDuration"].(time.Duration); ok && d1 > 0 {
		rateLimitDuration = d1
	} else if n1, ok := map1["rateLimitDuration"].(int); ok && n1 > 0 {
		rateLimitDuration = time.Duration(int64(n1)) * time.Second
	} else if n1, ok := map1["rateLimitDuration"].(int64); ok && n1 > 0 {
		rateLimitDuration = time.Duration(n1) * time.Second
	} else if s1, ok := map1["rateLimitDuration"].(string); ok && s1 != "" {
		rateLimitDuration = castx.ToDuration(s1)
	}

	return &RouteRule{
		httpMethod:        castx.ToString(map1["httpMethod"]),
		requestMapping:    castx.ToString(map1["requestMapping"]),
		regex:             castx.ToBool(map1["regex"]),
		pathVariableNames: pathVariableNames,
		jwtSettingsKey:    castx.ToString(map1["jwtSettingsKey"]),
		validateRules:     validateRules,
		failfast:          castx.ToBool(map1["failfast"]),
		rateLimitTotal:    castx.ToInt(map1["rateLimitTotal"], 0),
		rateLimitDuration: rateLimitDuration,
		rateLimitByIp:     castx.ToBool(map1["rateLimitByIp"]),
		customMiddlewares: make([]string, 0),
	}
}

func (rr *RouteRule) HttpMethod() string {
	return rr.httpMethod
}

func (rr *RouteRule) RequestMapping() string {
	return rr.requestMapping
}

func (rr *RouteRule) IsRegex() bool {
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

func (rr *RouteRule) RateLimitTotal() int {
	return rr.rateLimitTotal
}

func (rr *RouteRule) RateLimitDuration() time.Duration {
	return rr.rateLimitDuration
}

func (rr *RouteRule) RateLimitByIp() bool {
	return rr.rateLimitByIp
}

func (rr *RouteRule) WithCustomMiddleware(m string) {
	rr.customMiddlewares = append(rr.customMiddlewares, m)
}

func (rr *RouteRule) CustomMiddlewares() []string {
	return rr.customMiddlewares
}
