package httpx

import (
	"github.com/meiguonet/mgboot-go-common/util/castx"
)

var maxBodySize int64
var customJwtAuthExceptionHandler ExceptionHandler
var customValidateExceptionHandler ExceptionHandler

func WithMaxBodySize(arg0 interface{}) {
	var n1 int64

	if n2, ok := arg0.(int64); ok {
		n1 = n2
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		n1 = castx.ToDataSize(s1)
	}

	maxBodySize = n1
}

func MaxBodySize(args ...interface{}) int64 {
	if len(args) < 1 {
		n1 := maxBodySize

		if n1 < 1 {
			n1 = 8 * 1024 * 1024
		}

		return n1
	}

	var n1 int64

	if n2, ok := args[0].(int64); ok {
		n1 = n2
	} else if s1, ok := args[0].(string); ok && s1 != "" {
		n1 = castx.ToDataSize(s1)
	}

	maxBodySize = n1
	return 0
}

func WithJwtAuthExceptionHandler(handler ExceptionHandler) {
	customJwtAuthExceptionHandler = handler
}

func WithValidateExceptionHandler(handler ExceptionHandler)  {
	customValidateExceptionHandler = handler
}
