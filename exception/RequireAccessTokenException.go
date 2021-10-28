package BuiltinException

type RequireAccessTokenException struct {
	errorTips string
}

func NewRequireAccessTokenException(errorTips ...string) RequireAccessTokenException {
	var tips string

	if len(errorTips) > 0 {
		tips = errorTips[0]
	}

	if tips == "" {
		tips = "安全令牌缺失"
	}

	return RequireAccessTokenException{errorTips: tips}
}

func (ex RequireAccessTokenException) Error() string {
	return ex.errorTips
}
