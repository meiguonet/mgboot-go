package httpx

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/fsx"
	"github.com/meiguonet/mgboot-go-dal/ratelimiter"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"time"
)

type rateLimitMiddleware struct {
	total     int
	duration  time.Duration
	limitByIp bool
	luaFile   string
	cacheDir  string
}

func NewRateLimitMiddleware(settings map[string]interface{}) *rateLimitMiddleware {
	var total int

	if n1, ok := settings["total"].(int); ok && n1 > 0 {
		total = n1
	}

	var duration time.Duration

	if d1, ok := settings["duration"].(time.Duration); ok && d1 > 0 {
		duration = d1
	} else if n1, ok := settings["duration"].(int); ok && n1 > 0 {
		duration = time.Duration(int64(n1)) * time.Second
	} else if n1, ok := settings["duration"].(int64); ok && n1 > 0 {
		duration = time.Duration(n1) * time.Second
	} else if s1, ok := settings["duration"].(string); ok && s1 != "" {
		duration = castx.ToDuration(s1)
	}

	var limitByIp bool

	if b1, ok := settings["limitByIp"].(bool); ok {
		limitByIp = b1
	}

	var luaFile string

	if s1, ok := settings["luaFile"].(string); ok && s1 != "" {
		luaFile = fsx.GetRealpath(luaFile)
	}

	var cacheDir string

	if s1, ok := settings["cacheDir"].(string); ok && s1 != "" {
		cacheDir = fsx.GetRealpath(cacheDir)
	}

	return &rateLimitMiddleware{
		total:     total,
		duration:  duration,
		limitByIp: limitByIp,
		luaFile:   luaFile,
		cacheDir:  cacheDir,
	}
}

func (m *rateLimitMiddleware) GetName() string {
	return "builtin.RateLimit"
}

func (m *rateLimitMiddleware) GetType() int {
	return MiddlewareType.PreHandle
}

func (m *rateLimitMiddleware) GetOrder() int {
	return MiddlewareOrder.Highest
}

func (m *rateLimitMiddleware) PreHandle(req *Request, resp *Response) {
	if !req.Next() || resp.HasError() || m.total < 1 || m.duration < 1 {
		return
	}

	routeRule := req.GetRouteRule()

	if routeRule == nil {
		return
	}

	id := routeRule.HttpMethod() + "@" + routeRule.RequestMapping()

	if m.limitByIp {
		id += "@" + req.GetClientIp()
	}

	opts := ratelimiter.NewRatelimiterOptions(m.luaFile, m.cacheDir)
	result := ratelimiter.NewRatelimiter(id, m.total, m.duration, opts).GetLimit()
	remaining := castx.ToInt(result["remaining"])

	if remaining >= 0 {
		return
	}

	err := BuiltinException.NewRateLimitExceedException(map[string]interface{}{
		"total": castx.ToInt(result["total"]),
		"remaining": remaining,
		"retryAfter": castx.ToString(result["retryAfter"]),
	})

	resp.WithError(err)
	req.Next(false)
}

func (m *rateLimitMiddleware) PostHandle(_ *Request, _ *Response) {
}
