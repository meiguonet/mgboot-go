package mgboot

import (
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"github.com/meiguonet/mgboot-go/httpx"
	"net/http"
	"regexp"
	"strings"
)


var applicationJson = "application/json; charset=utf-8"
var textPlain = "text/plain; charset=utf-8"
var respBuf1 = []byte(`{"code":200}`)

func GetBuiltintExceptionName(name string) string {
	name = stringx.EnsureRight(name, "Exception")
	return stringx.EnsureLeft(name, "builtin.")
}

func Dispatch(w http.ResponseWriter, req *http.Request, modules []httpx.HandlerModule) {
	method := strings.ToUpper(req.Method)

	if method == "OPTIONS" {
		w.Header().Add("Content-Type", applicationJson)
		w.Write(respBuf1)
		return
	}

	requestUri := stringx.EnsureLeft(req.RequestURI, "/")
	handlers := make([]*httpx.HandlerEntry, 0)
	pathVariables := map[string]string{}

	for _, m := range modules {
		entries := m.GetHandlerEntries()

		for _, entry := range entries {
			routeRule := entry.GetRouteRule()

			if routeRule == nil {
				continue
			}

			if routeRule.RequestMapping() == requestUri {
				handlers = append(handlers, entry)
				continue
			}

			if routeRule.Regex() == "" {
				continue
			}

			re, err := regexp.Compile(routeRule.Regex())

			if err != nil {
				continue
			}

			matches := re.FindStringSubmatch(requestUri)

			if len(matches) < 2 {
				continue
			}

			handlers = append(handlers, entry)

			if len(pathVariables) > 0 {
				continue
			}

			pvNames := routeRule.PathVariableNames()

			if len(pvNames) != len(matches) - 1 {
				continue
			}

			for i, pvn := range pvNames {
				pathVariables[pvn] = matches[i + 1]
			}
		}
	}

	if len(handlers) < 1 {
		w.Header().Add("Content-Type", textPlain)
		w.WriteHeader(404)
		w.Write([]byte{})
		return
	}

	var handler *httpx.HandlerEntry

	for _, h := range handlers {
		if h.GetRouteRule().HttpMethod() == method {
			handler = h
			break
		}
	}

	if handler == nil {
		w.Header().Add("Content-Type", textPlain)
		w.WriteHeader(405)
		w.Write([]byte{})
		return
	}

	request := httpx.NewRequest(req).WithPathVariables(pathVariables)
	response := httpx.NewResponse(request, w)
	handler.HandleRequest(request, response)
}
