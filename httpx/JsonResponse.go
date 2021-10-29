package httpx

import (
	"github.com/meiguonet/mgboot-go-common/util/jsonx"
	"strings"
)

type JsonResponse struct {
	payload interface{}
}

func NewJsonResponse(payload interface{}) JsonResponse {
	return JsonResponse{payload: payload}
}

func (p JsonResponse) GetContentType() string {
	return "application/json; charset=utf-8"
}

func (p JsonResponse) GetContents() interface{} {
	if s1, ok := p.payload.(string); ok && p.isJson(s1) {
		return s1
	}

	opts := jsonx.NewToJsonOption().HandleTimeField().StripZeroTimePart()
	contents := jsonx.ToJson(p.payload, opts)

	if p.isJson(contents) {
		return contents
	}

	return "{}"
}

func (p JsonResponse) isJson(contents string) bool {
	var flag bool

	if strings.HasPrefix(contents, "{") && strings.HasPrefix(contents, "}") {
		flag = true
	} else if strings.HasPrefix(contents, "[") && strings.HasPrefix(contents, "]") {
		flag = true
	}

	return flag
}
