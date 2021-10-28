package BuiltinException

type AccessTokenInvalidException struct {
	errorTips string
}

func NewAccessTokenInvalidException(errorTips ...string) AccessTokenInvalidException {
	var tips string

	if len(errorTips) > 0 {
		tips = errorTips[0]
	}

	if tips == "" {
		tips = "不是有效安全令牌"
	}

	return AccessTokenInvalidException{errorTips: tips}
}

func (ex AccessTokenInvalidException) Error() string {
	return ex.errorTips
}
