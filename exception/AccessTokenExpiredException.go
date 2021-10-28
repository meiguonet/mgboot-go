package BuiltinException

type AccessTokenExpiredException struct {
	errorTips string
}

func NewAccessTokenExpiredException(errorTips ...string) AccessTokenExpiredException {
	var tips string

	if len(errorTips) > 0 {
		tips = errorTips[0]
	}

	if tips == "" {
		tips = "安全令牌已失效"
	}

	return AccessTokenExpiredException{errorTips: tips}
}

func (ex AccessTokenExpiredException) Error() string {
	return ex.errorTips
}
