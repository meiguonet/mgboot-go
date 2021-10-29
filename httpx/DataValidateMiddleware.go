package httpx

import (
	"github.com/meiguonet/mgboot-go-common/util/validatex"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareOrder"
	"github.com/meiguonet/mgboot-go/enum/MiddlewareType"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
)

type dataValidateMiddleware struct {
}

func NewDataValidateMiddleware() *dataValidateMiddleware {
	return &dataValidateMiddleware{}
}

func (m *dataValidateMiddleware) GetName() string {
	return "builtin.DataValidate"
}

func (m *dataValidateMiddleware) GetType() int {
	return MiddlewareType.PreHandle
}

func (m *dataValidateMiddleware) GetOrder() int {
	return MiddlewareOrder.Highest
}

func (m *dataValidateMiddleware) PreHandle(req *Request, resp *Response) {
	if !req.Next() || resp.HasError() {
		return
	}

	routeRule := req.GetRouteRule()

	if routeRule == nil {
		return
	}

	if len(routeRule.ValidateRules()) < 1 {
		return
	}

	validator := validatex.NewValidator()
	data := req.GetMap()

	if routeRule.IsFailfast() {
		errorTips := validatex.FailfastValidate(validator, data, routeRule.ValidateRules())
		err := BuiltinException.NewValidateException(errorTips, true)

		if errorTips != "" {
			req.Next(false)
			resp.WithError(err)
		}

		return
	}

	validateErrors := validatex.Validate(validator, data, routeRule.ValidateRules())

	if len(validateErrors) > 0 {
		err := BuiltinException.NewValidateException(validateErrors)
		req.Next(false)
		resp.WithError(err)
	}
}

func (m *dataValidateMiddleware) PostHandle(_ *Request, _ *Response) {
}
