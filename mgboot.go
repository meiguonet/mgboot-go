package mgboot

import (
	"github.com/meiguonet/mgboot-go-common/logx"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go/httpx"
	"github.com/meiguonet/mgboot-go/securityx"
	"os"
)

var maxBodySize int64
var corsSettings *securityx.CorsSettings
var jwtPublicKeyPemFile string
var jwtPrivateKeyPemFile string
var jwtSettings = map[string]*securityx.JwtSettings{}
var middlewares = make([]httpx.Middleware, 0)
var exceptionHandlers = make([]httpx.ExceptionHandler, 0)
var requestLogLogger logx.Logger
var executeTimeLogLogger logx.Logger

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
	entries := []httpx.Middleware{
		httpx.NewRequestLogMiddleware(),
		httpx.NewExecuteTimeLogMiddleware(),
	}

	middlewares = append(middlewares, entries...)
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

func Middlewares() []httpx.Middleware {
	return middlewares
}

func WithExceptionHandler(handler httpx.ExceptionHandler) {
	exceptionHandlers = append(exceptionHandlers, handler)
}

func WithExceptionHandlers(handlers []httpx.ExceptionHandler) {
	if len(handlers) < 1 {
		return
	}

	exceptionHandlers = append(exceptionHandlers, handlers...)
}

func ExceptionHandlers() []httpx.ExceptionHandler {
	return exceptionHandlers
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
