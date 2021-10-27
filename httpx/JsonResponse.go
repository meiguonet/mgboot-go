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

func (p JsonResponse) GetContents() string {
	opts := jsonx.ToJsonOption{
		HandleTimeField:   true,
		StripZeroTimePart: true,
	}

	contents := jsonx.ToJson(p.payload, opts)
	var isJson bool

	if strings.HasPrefix(contents, "{") && strings.HasPrefix(contents, "}") {
		isJson = true
	} else if strings.HasPrefix(contents, "[") && strings.HasPrefix(contents, "]") {
		isJson = true
	}

	if !isJson {
		return "{}"
	}

	return contents
}
