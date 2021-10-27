package mvc

import "github.com/meiguonet/mgboot-go-common/util/castx"

type HandlerFuncArgInfo struct {
	name             string
	typ              string
	request          bool
	jwt              bool
	clientIp         bool
	httpHeaderName   string
	jwtClaimName     string
	pathVariableName string
	requestParamName string
	securityMode     int
	paramMapRules    []string
	formFieldName    string
	needRequestBody  bool
}

func NewHandlerFuncArgInfo(map1 map[string]interface{}) HandlerFuncArgInfo {
	paramMapRules := make([]string, 0)

	if a1, ok := map1["paramMapRules"].([]string); ok && len(a1) > 0 {
		paramMapRules = a1
	}

	return HandlerFuncArgInfo{
		name:             castx.ToString(map1["name"]),
		typ:              castx.ToString(map1["type"]),
		request:          castx.ToBool(map1["request"]),
		jwt:              castx.ToBool(map1["jwt"]),
		clientIp:         castx.ToBool(map1["clientIp"]),
		httpHeaderName:   castx.ToString(map1["httpHeaderName"]),
		jwtClaimName:     castx.ToString(map1["jwtClaimName"]),
		pathVariableName: castx.ToString(map1["pathVariableName"]),
		requestParamName: castx.ToString(map1["requestParamName"]),
		securityMode:     castx.ToInt(map1["securityMode"], 0),
		paramMapRules:    paramMapRules,
		formFieldName:    castx.ToString(map1["formFieldName"]),
		needRequestBody:  castx.ToBool(map1["needRequestBody"]),
	}
}

func (info HandlerFuncArgInfo) Name() string {
	return info.name
}

func (info HandlerFuncArgInfo) Type() string {
	return info.typ
}

func (info HandlerFuncArgInfo) IsRequest() bool {
	return info.request
}

func (info HandlerFuncArgInfo) IsJwt() bool {
	return info.jwt
}

func (info HandlerFuncArgInfo) IsClientIp() bool {
	return info.clientIp
}

func (info HandlerFuncArgInfo) HttpHeaderName() string {
	return info.httpHeaderName
}

func (info HandlerFuncArgInfo) JwtClaimName() string {
	return info.jwtClaimName
}

func (info HandlerFuncArgInfo) PathVariableName() string {
	return info.pathVariableName
}

func (info HandlerFuncArgInfo) RequestParamName() string {
	return info.requestParamName
}

func (info HandlerFuncArgInfo) SecurityMode() int {
	return info.securityMode
}

func (info HandlerFuncArgInfo) ParamMapRules() []string {
	return info.paramMapRules
}

func (info HandlerFuncArgInfo) FormFieldName() string {
	return info.formFieldName
}

func (info HandlerFuncArgInfo) IsNeedRequestBody() bool {
	return info.needRequestBody
}
