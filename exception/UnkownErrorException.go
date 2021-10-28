package BuiltinException

type UnkownErrorException struct {
	errorTips string
}

func NewUnkownErrorException(errorTips ...string) UnkownErrorException {
	var tips string

	if len(errorTips) > 0 {
		tips = errorTips[0]
	}

	if tips == "" {
		tips = "unkown error found"
	}

	return UnkownErrorException{errorTips: tips}
}

func (ex UnkownErrorException) Error() string {
	return ex.errorTips
}
