package httpx

import (
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/securityx"
)

type jwtAuthMiddleware struct {
	settings *securityx.JwtSettings
}

func NewJwtAuthMiddleware(settings *securityx.JwtSettings) *jwtAuthMiddleware {
	return &jwtAuthMiddleware{settings: settings}
}

func (m *jwtAuthMiddleware) GetName() string {
	return "builtin.JwtAuth"
}

func (m *jwtAuthMiddleware) GetType() int {
	return MiddlewareType.PreHandle
}

func (m *jwtAuthMiddleware) GetOrder() int {
	return MiddlewareOrder.Highest
}

func (m *jwtAuthMiddleware) PreHandle(req *Request, resp *Response) {
	if !req.Next() || resp.HasError() {
		return
	}

	routeRule := req.GetRouteRule()

	if routeRule == nil {
		return
	}

	key := routeRule.JwtSettingsKey()

	if key == "" {
		return
	}

	if m.settings == nil {
		return
	}

	token := req.GetJwt()

	if token == nil {
		req.Next(false)
		resp.WithError(BuiltinException.NewRequireAccessTokenException())
		return
	}

	if err := securityx.VerifyJsonWebToken(token, m.settings); err != nil {
		req.Next(false)
		resp.WithError(err)
	}
}

func (m *jwtAuthMiddleware) PostHandle(_ *Request, _ *Response) {
}
