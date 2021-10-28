package BuiltinException

type ValidateException struct {
	errorTips      string
	validateErrors map[string]string
	failfast       bool
}

func NewValidateException(args ...interface{}) ValidateException {
	var errorTips string

	for _, arg := range args {
		if s1, ok := arg.(string); ok && s1 != "" {
			errorTips = s1
			break
		}
	}

	if errorTips == "" {
		errorTips = "数据完整性验证错误"
	}

	validateErrors := map[string]string{}

	for _, arg := range args {
		if map1, ok := arg.(map[string]string); ok && len(map1) > 0 {
			validateErrors = map1
			break
		}
	}

	var failfast bool

	for _, arg := range args {
		if b1, ok := arg.(bool); ok {
			failfast = b1
			break
		}
	}

	return ValidateException{
		errorTips:      errorTips,
		validateErrors: validateErrors,
		failfast:       failfast,
	}
}

func (ex ValidateException) Error() string {
	return ex.errorTips
}

func (ex ValidateException) GetValidateErrors() map[string]string {
	return ex.validateErrors
}

func (ex ValidateException) IsFailfast() bool {
	return ex.failfast
}
