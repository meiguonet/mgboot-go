package mgboot

import (
	"github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"github.com/meiguonet/mgboot-go/httpx"
	BuiltinExceptionHandler "github.com/meiguonet/mgboot-go/httpx/ExceptionHandler"
	BuiltinMiddleware "github.com/meiguonet/mgboot-go/httpx/middleware"
	"github.com/meiguonet/mgboot-go/securityx"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var maxBodySize int64
var corsSettings *securityx.CorsSettings
var jwtPublicKeyPemFile string
var jwtPrivateKeyPemFile string
var jwtSettings = map[string]*securityx.JwtSettings{}
var middlewares = make([]httpx.Middleware, 0)
var exceptionHandlers = make([]httpx.ExceptionHandler, 0)
var runtimeLogger logx.Logger
var requestLogLogger logx.Logger
var executeTimeLogLogger logx.Logger
var applicationJson = "application/json; charset=utf-8"
var textPlain = "text/plain; charset=utf-8"
var response1 = []byte(`{"code":200}`)

func WithMaxBodySize(arg0 interface{}) {
	var n1 int64

	if n2, ok := arg0.(int64); ok {
		n1 = n2
	} else if s1, ok := arg0.(string); ok && s1 != "" {
		n1 = castx.ToDataSize(s1)
	}

	maxBodySize = n1
}

func MaxBodySize() int64 {
	n1 := maxBodySize

	if n1 < 1 {
		n1 = 8 * 1024 * 1024
	}

	return n1
}

func WithCorsSettings(arg0 interface{}) {
	if settings, ok := arg0.(*securityx.CorsSettings); ok && settings != nil {
		corsSettings = settings
	} else if settings, ok := arg0.(map[string]interface{}); ok && len(settings) > 0 {
		corsSettings = securityx.NewCorsSettings(settings)
	}
}

func CorsSettings() *securityx.CorsSettings {
	return corsSettings
}

func WithJwtPublicKeyPemFile(fpath string) {
	if stat, err := os.Stat(fpath); err != nil || stat.IsDir() {
		return
	}

	jwtPublicKeyPemFile = fpath
}

func JwtPublicKeyPemFile() string {
	return jwtPublicKeyPemFile
}

func WithJwtPrivateKeyPemFile(fpath string) {
	if stat, err := os.Stat(fpath); err != nil || stat.IsDir() {
		return
	}

	jwtPrivateKeyPemFile = fpath
}

func JwtPrivateKeyPemFile() string {
	return jwtPrivateKeyPemFile
}

func WithJwtSettings(key string, settings interface{}) {
	if st, ok := settings.(*securityx.JwtSettings); ok && st != nil {
		jwtSettings[st.Key()] = st
	} else if map1, ok := settings.(map[string]interface{}); ok && len(map1) > 0 {
		jwtSettings[key] = securityx.NewJwtSettings(key, map1)
	}
}

func JwtSettings(key string) *securityx.JwtSettings {
	return jwtSettings[key]
}

func WithBuiltinMiddlewares() {
	middlewares = []httpx.Middleware{
		BuiltinMiddleware.NewJwtAuthMiddleware(),
		BuiltinMiddleware.NewDataValidateMiddleware(),
		BuiltinMiddleware.NewRequestLogMiddleware(),
		BuiltinMiddleware.NewExecuteTimeLogMiddleware(),
	}
}

func WithMiddleware(m httpx.Middleware) {
	if len(middlewares) < 1 {
		WithBuiltinMiddlewares()
	}

	middlewares = append(middlewares, m)
}

func WithMiddlewares(entries []httpx.Middleware) {
	if len(middlewares) < 1 {
		WithBuiltinMiddlewares()
	}

	if len(entries) < 1 {
		return
	}

	middlewares = append(middlewares, entries...)
}

func ReplaceBuiltinValidateMiddleware(m httpx.Middleware) {
	entries := make([]httpx.Middleware, 0)

	for _, middleware := range middlewares {
		if m.GetName() == "builtin.DataValidate" {
			entries = append(entries, m)
		} else {
			entries = append(entries, middleware)
		}
	}

	middlewares = entries
}

func Middlewares() []httpx.Middleware {
	return middlewares
}

func WithBuiltinExceptionHandlers() {
	exceptionHandlers = []httpx.ExceptionHandler{
		BuiltinExceptionHandler.NewAccessTokenExpiredHandler(),
		BuiltinExceptionHandler.NewAccessTokenInvalidHandler(),
		BuiltinExceptionHandler.NewDbExceptionHandler(),
		BuiltinExceptionHandler.NewHttpErrorHandler(),
		BuiltinExceptionHandler.NewRequireAccessTokenHandler(),
		BuiltinExceptionHandler.NewUnkownErrorHandler(),
		BuiltinExceptionHandler.NewValidateHandler(),
	}
}

func WithExceptionHandler(handler httpx.ExceptionHandler) {
	if len(exceptionHandlers) < 1 {
		WithBuiltinExceptionHandlers()
	}

	idx := -1

	for n1, h := range exceptionHandlers {
		if h.GetExceptionName() == handler.GetExceptionName() {
			idx = n1
		}
	}

	if idx < 0 {
		exceptionHandlers = append(exceptionHandlers, handler)
		return
	}

	handlers := make([]httpx.ExceptionHandler, 0)

	for n1, h := range exceptionHandlers {
		if n1 == idx {
			handlers = append(handlers, handler)
		} else {
			handlers = append(handlers, h)
		}
	}

	exceptionHandlers = handlers
}

func WithExceptionHandlers(handlers []httpx.ExceptionHandler) {
	for _, handler := range handlers {
		WithExceptionHandler(handler)
	}
}

func ExceptionHandlers() []httpx.ExceptionHandler {
	return exceptionHandlers
}

func WithRuntimeLogLogger(logger logx.Logger) {
	runtimeLogger = logger
}

func RuntimeLogger() logx.Logger {
	return runtimeLogger
}

func WithRequestLogLogger(logger logx.Logger) {
	requestLogLogger = logger
}

func RequestLogLogger() logx.Logger {
	return requestLogLogger
}

func WithExecuteTimeLogLogger(logger logx.Logger) {
	executeTimeLogLogger = logger
}

func ExecuteTimeLogLogger() logx.Logger {
	return executeTimeLogLogger
}

func GetBuiltintExceptionName(name string) string {
	name = stringx.EnsureRight(name, "Exception")
	return stringx.EnsureLeft(name, "BuiltinException.")
}

func Dispatch(w http.ResponseWriter, req *http.Request, modules []httpx.HandlerModule) {
	method := strings.ToUpper(req.Method)

	if method == "OPTIONS" {
		w.Header().Add("Content-Type", applicationJson)
		w.Write(response1)
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

	handler.HandleRequest(w, req, pathVariables)
}
