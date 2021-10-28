package BuiltinMiddleware

import (
	"github.com/meiguonet/mgboot-go"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/httpx"
	"github.com/meiguonet/mgboot-go/securityx"
)

type jwtAuthMiddleware struct {
}

func NewJwtAuthMiddleware() *jwtAuthMiddleware {
	return &jwtAuthMiddleware{}
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

func (m *jwtAuthMiddleware) PreHandle(req *httpx.Request, resp *httpx.Response) {
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

	settings := mgboot.JwtSettings(key)

	if settings == nil {
		return
	}

	token := req.GetJwt()

	if token == nil {
		req.Next(false)
		resp.WithError(BuiltinException.NewRequireAccessTokenException())
		return
	}

	if err := securityx.VerifyJsonWebToken(token, settings); err != nil {
		req.Next(false)
		resp.WithError(err)
	}
}

func (m *jwtAuthMiddleware) PostHandle(_ *httpx.Request, _ *httpx.Response) {
}
