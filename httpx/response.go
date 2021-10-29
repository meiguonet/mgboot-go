package httpx

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/slicex"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	BuiltinException "github.com/meiguonet/mgboot-go/exception"
	"github.com/meiguonet/mgboot-go/logx"
	"github.com/meiguonet/mgboot-go/securityx"
	"net/http"
	"strings"
)

var payloadNilResponse = []byte(`{"code":500,"msg":"response payload is nil"}`)
var unsupportedPayloadContentsResponse = []byte(`{"code":500,"msg":"unsupported response payload contents"}`)
var emptyResponse = make([]byte, 0)

type Response struct {
	request           *Request
	out               http.ResponseWriter
	payload           ResponsePayload
	err               error
	extraHeaders      map[string]string
	exceptionHandlers []ExceptionHandler
	corsSettings      *securityx.CorsSettings
}

func NewResponse(request *Request, out http.ResponseWriter) *Response {
	exceptionHandlers := []ExceptionHandler{
		NewDbExceptionHandler(),
		NewHttpErrorHandler(),
		NewUnkownErrorHandler(),
	}

	if customJwtAuthExceptionHandler != nil {
		exceptionHandlers = append(exceptionHandlers, customJwtAuthExceptionHandler)
	} else {
		exceptionHandlers = append(exceptionHandlers, NewAccessTokenExpiredHandler())
		exceptionHandlers = append(exceptionHandlers, NewAccessTokenInvalidHandler())
		exceptionHandlers = append(exceptionHandlers, NewRequireAccessTokenHandler())
	}

	if customValidateExceptionHandler != nil {
		exceptionHandlers = append(exceptionHandlers, customValidateExceptionHandler)
	} else {
		exceptionHandlers = append(exceptionHandlers, NewValidateExceptionHandler())
	}

	resp := &Response{
		request:           request,
		out:               out,
		extraHeaders:      map[string]string{},
		exceptionHandlers: exceptionHandlers,
	}

	return resp
}

func (resp *Response) WithPayload(payload ResponsePayload) *Response {
	resp.payload = payload
	return resp
}

func (resp *Response) WithError(err error) *Response {
	resp.err = err
	return resp
}

func (resp *Response) WithExtraHeader(headerName, headerValue string) *Response {
	if resp.extraHeaders == nil {
		resp.extraHeaders = map[string]string{}
	}

	headerName = stringx.Ucwords(headerName, "-", "-")
	resp.extraHeaders[headerName] = headerValue
	return resp
}

func (resp *Response) WithExtraHeaders(headers map[string]string) *Response {
	if len(headers) > 0 {
		for headerName, headerValue := range headers {
			resp.WithExtraHeader(headerName, headerValue)
		}
	}

	return resp
}

func (resp *Response) WithExceptionHandler(handler ExceptionHandler) *Response {
	resp.exceptionHandlers = append(resp.exceptionHandlers, handler)
	return resp
}

func (resp *Response) WithExceptionHandlers(handlers []ExceptionHandler) *Response {
	if len(handlers) > 0 {
		resp.exceptionHandlers = append(resp.exceptionHandlers, handlers...)
	}

	return resp
}

func (resp *Response) WithCorsSettings(corsSettings *securityx.CorsSettings) *Response {
	if corsSettings != nil {
		resp.corsSettings = corsSettings
	}

	return resp
}

func (resp *Response) HasError() bool {
	return resp.err != nil
}

func (resp *Response) Send() {
	if resp.corsSettings != nil {
		resp.addCorsSupport()
	}

	if resp.err != nil {
		if ex, ok := resp.err.(BuiltinException.RateLimitExceedException); ok {
			resp.WithExtraHeader("X-Ratelimit-Limit", fmt.Sprintf("%d", ex.Total()))
			resp.WithExtraHeader("X-Ratelimit-Remaining", fmt.Sprintf("%d", ex.Remaining()))

			if ex.RetryAfter() != "" {
				resp.WithExtraHeader("Retry-After", ex.RetryAfter())
			}

			resp.sendWithHttpErrorCode(429)
			return
		}

		var handler ExceptionHandler

		if len(resp.exceptionHandlers) > 0 {
			for _, h := range resp.exceptionHandlers {
				if h.MatchException(resp.err) {
					handler = h
					break
				}
			}
		}

		if handler == nil {
			resp.writeErrorLog(resp.err)
			resp.err = BuiltinException.NewUnkownErrorException()
			resp.Send()
			return
		}

		resp.payload = handler.HandleException(resp.err)
		resp.err = nil
		resp.Send()
		return
	}

	if resp.payload == nil {
		resp.sendWithPayloadNilError()
		return
	}

	contents := resp.payload.GetContents()

	if errCode, ok := contents.(int); ok {
		resp.sendWithHttpErrorCode(errCode)
		return
	}

	if payload, ok := contents.(ImageResponse); ok {
		resp.sendImage(payload)
		return
	}

	if payload, ok := contents.(AttachmentResponse); ok {
		resp.sendAttachment(payload)
		return
	}

	if payload, ok := contents.(HttpError); ok {
		resp.sendWithHttpErrorCode(payload.GetStatusCode())
		return
	}

	s1, ok := contents.(string)

	if !ok {
		resp.sendWithUnsupportedPayloadContentsError()
		return
	}

	resp.sendString(resp.payload.GetContentType(), s1)
}

func (resp *Response) NeedCorsSupport() bool {
	methods := []string{"PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}

	if slicex.InStringSlice(resp.request.GetMethod(), methods) {
		return true
	}

	contentType := strings.ToLower(resp.request.GetHeader("Content-Type"))

	if strings.Contains(contentType, "application/x-www-form-urlencoded") ||
		strings.Contains(contentType, "multipart/form-data") ||
		strings.Contains(contentType, "text/plain") {
		return true
	}

	headerNames := []string{
		"Accept",
		"Accept-Language",
		"Content-Language",
		"DPR",
		"Downlink",
		"Save-Data",
		"Viewport-Widt",
		"Width",
	}

	for headerName := range resp.request.GetHeaders() {
		if slicex.InStringSlice(headerName, headerNames) {
			return true
		}
	}

	return false
}

func (resp *Response) addCorsSupport() {
	settings := resp.corsSettings

	if settings == nil {
		return
	}

	allowedOrigins := settings.AllowedOrigins()

	if slicex.InStringSlice("*", allowedOrigins) {
		resp.out.Header().Add("Access-Control-Allow-Origin", "*")
	} else {
		resp.out.Header().Add("Access-Control-Allow-Origin", strings.Join(allowedOrigins, ", "))
	}

	allowedHeaders := settings.AllowedHeaders()

	if len(allowedHeaders) > 0 {
		resp.out.Header().Add("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	}

	exposedHeaders := settings.ExposedHeaders()

	if len(exposedHeaders) > 0 {
		resp.out.Header().Add("Access-Control-Expose-Headers", strings.Join(exposedHeaders, ", "))
	}

	maxAge := settings.MaxAge()

	if maxAge > 0 {
		n1 := castx.ToInt64(maxAge.Seconds())
		resp.out.Header().Add("Access-Control-Max-Age", fmt.Sprintf("%d", n1))
	}

	if settings.AllowCredentials() {
		resp.out.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}

func (resp *Response) sendWithPayloadNilError() {
	resp.WithExtraHeader("Content-Type", "application/json; charset=utf-8")

	for headerName, headerValue := range resp.extraHeaders {
		resp.out.Header().Add(headerName, headerValue)
	}

	resp.out.Write(payloadNilResponse)
}

func (resp *Response) sendWithUnsupportedPayloadContentsError() {
	resp.WithExtraHeader("Content-Type", "application/json; charset=utf-8")

	for headerName, headerValue := range resp.extraHeaders {
		resp.out.Header().Add(headerName, headerValue)
	}

	resp.out.Write(unsupportedPayloadContentsResponse)
}

func (resp *Response) sendWithHttpErrorCode(code int) {
	resp.WithExtraHeader("Content-Type", "text/plain")

	for headerName, headerValue := range resp.extraHeaders {
		resp.out.Header().Add(headerName, headerValue)
	}

	resp.out.WriteHeader(code)
	resp.out.Write(emptyResponse)
}

func (resp *Response) sendImage(payload ImageResponse) {
	if len(payload.Buffer()) < 1 || payload.GetContentType() == "" {
		resp.sendWithHttpErrorCode(400)
		return
	}

	resp.WithExtraHeader("Content-Type", payload.GetContentType())

	for headerName, headerValue := range resp.extraHeaders {
		resp.out.Header().Add(headerName, headerValue)
	}

	resp.out.Write(payload.Buffer())
}

func (resp *Response) sendAttachment(payload AttachmentResponse) {
	if len(payload.Buffer()) < 1 || payload.AttachmentFileName() == "" {
		resp.sendWithHttpErrorCode(400)
		return
	}

	if resp.corsSettings != nil {
		resp.addCorsSupport()
	}

	disposition := fmt.Sprintf(`attachment; filename="%s"`, payload.AttachmentFileName())
	resp.WithExtraHeader("Content-Type", payload.GetContentType())
	resp.WithExtraHeader("Content-Length", fmt.Sprintf("%d", len(payload.Buffer())))
	resp.WithExtraHeader("Content-Transfer-Encoding", "binary")
	resp.WithExtraHeader("Content-Disposition", disposition)
	resp.WithExtraHeader("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	resp.WithExtraHeader("Pragma", "public")

	for headerName, headerValue := range resp.extraHeaders {
		resp.out.Header().Add(headerName, headerValue)
	}

	resp.out.Write(payload.Buffer())
}

func (resp *Response) sendString(contentType, contents string) {
	if contentType == "" {
		contentType = "text/plain; charset=utf-8"
	}

	resp.WithExtraHeader("Content-Type", contentType)

	for headerName, headerValue := range resp.extraHeaders {
		resp.out.Header().Add(headerName, headerValue)
	}

	resp.out.Write([]byte(contents))
}

func (resp *Response) writeErrorLog(arg0 interface{}) {
	var msg string

	if s1, ok := arg0.(string); ok {
		msg = s1
	} else if err, ok := arg0.(error); ok {
		msg = resp.getStacktrace(err)
	}

	if msg == "" {
		return
	}

	logx.Error(msg)
}

func (resp *Response) getStacktrace(arg0 interface{}) string {
	if arg0 == nil {
		return ""
	}

	var stacktrace string

	switch t := arg0.(type) {
	case *errors.Error:
		stacktrace = t.ErrorStack()
	case error:
		stacktrace = errors.New(t).ErrorStack()
	}

	if stacktrace == "" {
		return ""
	}

	stacktrace = strings.ReplaceAll(stacktrace, "\r", "")
	lines := strings.Split(stacktrace, "\n")
	var trace []string

	if strings.Contains(stacktrace, "src/runtime/panic.go") {
		n1 := -1

		for i := 0; i < len(lines); i++ {
			if i == 0 {
				trace = append(trace, lines[i])
				continue
			}

			if strings.Contains(lines[i], "src/runtime/panic.go") {
				n1 = i
				continue
			}

			if strings.Contains(lines[i], "src/runtime/proc.go") ||
				strings.Contains(lines[i], "src/runtime/asm_amd64") {
				break
			}

			if n1 < 0 || i <= n1 + 1 {
				continue
			}

			trace = append(trace, lines[i])
		}
	} else {
		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], "src/runtime/proc.go") ||
				strings.Contains(lines[i], "src/runtime/asm_amd64") {
				break
			}

			trace = append(trace, lines[i])
		}
	}

	if len(trace) < 1 {
		return ""
	}

	return strings.Join(trace, "\n")
}
