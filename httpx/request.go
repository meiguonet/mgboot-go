package httpx

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/meiguonet/mgboot-go"
	"github.com/meiguonet/mgboot-go-common/enum/RegexConst"
	"github.com/meiguonet/mgboot-go-common/util/castx"
	"github.com/meiguonet/mgboot-go-common/util/jsonx"
	"github.com/meiguonet/mgboot-go-common/util/mapx"
	"github.com/meiguonet/mgboot-go-common/util/numberx"
	"github.com/meiguonet/mgboot-go-common/util/slicex"
	"github.com/meiguonet/mgboot-go-common/util/stringx"
	"github.com/meiguonet/mgboot-go/mvc"
	"github.com/meiguonet/mgboot-go/securityx"
	"math"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Request struct {
	rawRequest    *http.Request
	method        string
	headers       map[string]string
	queryParams   map[string]string
	formData      map[string]string
	formFiles     map[string]*multipart.FileHeader
	pathVariables map[string]string
	rawBody       []byte
	routeRule     *mvc.RouteRule
	middlewares   []Middleware
	execStart     time.Time
	next          bool
}

func NewRequest(req *http.Request) *Request {
	method := strings.ToUpper(req.Method)
	headers := map[string]string{}

	for headerName, headerValues := range req.Header {
		if len(headerValues) < 1 {
			continue
		}

		headerName = stringx.Ucwords(headerName, "-", "-")
		headers[headerName] = strings.Join(headerValues, "; ")
	}

	queryParams := map[string]string{}

	for name, values := range req.URL.Query() {
		if len(values) < 1 {
			continue
		}

		queryParams[name] = values[0]
	}

	formData := map[string]string{}
	formFiles := map[string]*multipart.FileHeader{}

	if method == "POST" {
		for name, values := range req.PostForm {
			if len(values) < 1 {
				continue
			}

			formData[name] = values[0]
		}

		if strings.Contains(strings.ToLower(headers["Content-Type"]), "multipart/form-data") {
			req.ParseMultipartForm(mgboot.MaxBodySize())

			for name, entries := range req.MultipartForm.File {
				if len(entries) < 0 {
					continue
				}

				formFiles[name] = entries[0]
			}
		}
	}

	return &Request{
		rawRequest:    req,
		method:        method,
		headers:       headers,
		queryParams:   queryParams,
		formData:      formData,
		formFiles:     formFiles,
		pathVariables: map[string]string{},
		rawBody:       buildRawBody(method, headers["Content-Type"], formData, req),
		middlewares:   mgboot.Middlewares(),
		execStart:     time.Now(),
		next:          true,
	}
}

func buildRawBody(method, contentType string, formData map[string]string, req *http.Request) []byte {
	methods := []string{"POST", "PUT", "PATCH", "DELETE"}

	if !slicex.InStringSlice(method, methods) {
		return make([]byte, 0)
	}

	contentType = strings.ToLower(contentType)
	isPostForm := strings.Contains(contentType, "application/x-www-form-urlencoded")
	isMultipartForm := strings.Contains(contentType, "multipart/form-data")

	if method == "POST" && (isPostForm || isMultipartForm) {
		sb := make([]string, 0)

		for name, value := range formData {
			sb = append(sb, fmt.Sprintf("%s=%s", name, value))
		}

		return []byte(strings.Join(sb, "&"))
	}

	isJson := strings.Contains(contentType, "application/json")
	isXml1 := strings.Contains(contentType, "application/xml")
	isXml2 := strings.Contains(contentType, "text/xml")

	if !isJson && !isXml1 && !isXml2 {
		return make([]byte, 0)
	}

	body := make([]byte, 0)

	for {
		buf := make([]byte, 0, 1024)
		n1, err := req.Body.Read(buf)

		if err != nil || n1 < 1 {
			break
		}

		body = append(body, buf[:n1]...)
	}

	req.Body.Close()
	return body
}

func (r *Request) WithPathVariables(map1 map[string]string) *Request {
	if len(map1) > 0 {
		r.pathVariables = map1
	}

	return r
}

func (r *Request) WithRouteRule(rr *mvc.RouteRule) *Request {
	if rr != nil {
		r.routeRule = rr
	}

	return r
}

func (r *Request) GetMethod() string {
	return r.method
}

func (r *Request) GetHeaders() map[string]string {
	return r.headers
}

func (r *Request) GetHeader(headerName string) string {
	headerName = strings.ToLower(headerName)

	for key, headerValue := range r.headers {
		if strings.ToLower(key) == headerName {
			return headerValue
		}
	}

	return ""
}

func (r *Request) GetRequestUrl(withQueryString ...bool) string {
	var b1 bool

	if len(withQueryString) > 0 {
		b1 = withQueryString[0]
	}

	s1 := stringx.EnsureLeft(r.rawRequest.RequestURI, "/")

	if !b1 {
		return s1
	}

	s2 := r.GetQueryString()

	if s2 != "" {
		return s1 + "?" + s2
	}

	return s1
}

func (r *Request) GetQueryString(urlencode ...bool) string {
	if len(r.queryParams) < 1 {
		return ""
	}

	var b1 bool

	if len(urlencode) > 0 {
		b1 = urlencode[0]
	}

	if !b1 {
		sb := make([]string, 0)

		for name, value := range r.queryParams {
			sb = append(sb, fmt.Sprintf("%s=%s", name, value))
		}

		return strings.Join(sb, "&")
	}

	values := url.Values{}

	for name, value := range r.queryParams {
		values[name] = []string{value}
	}

	return values.Encode()
}

func (r *Request) GetClientIp() string {
	ip := r.GetHeader("X-Forwarded-For")

	if ip == "" {
		ip = r.GetHeader("X-Real-IP")
	}

	if ip == "" {
		ip = r.rawRequest.RemoteAddr
	}

	regex1 := regexp.MustCompile(RegexConst.CommaSep)
	parts := regex1.Split(strings.TrimSpace(ip), -1)

	if len(parts) < 1 {
		return ""
	}

	return strings.TrimSpace(parts[0])
}

func (r *Request) GetInt(name string, defaultValue ...int) int {
	n1 := math.MinInt32

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.queryParams[name]; ok && s1 != "" {
		return castx.ToInt(s1, n1)
	}

	if s1, ok := r.formData[name]; ok && s1 != "" {
		return castx.ToInt(s1, n1)
	}

	return n1
}

func (r *Request) GetInt64(name string, defaultValue ...int64) int64 {
	n1 := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.queryParams[name]; ok && s1 != "" {
		return castx.ToInt64(s1, n1)
	}

	if s1, ok := r.formData[name]; ok && s1 != "" {
		return castx.ToInt64(s1, n1)
	}

	return n1
}

func (r *Request) GetFloat32(name string, defaultValue ...float32) float32 {
	n1 := castx.ToFloat32(math.SmallestNonzeroFloat32)

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.queryParams[name]; ok && s1 != "" {
		return castx.ToFloat32(s1, n1)
	}

	if s1, ok := r.formData[name]; ok && s1 != "" {
		return castx.ToFloat32(s1, n1)
	}

	return n1
}

func (r *Request) GetFloat64(name string, defaultValue ...float64) float64 {
	n1 := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.queryParams[name]; ok && s1 != "" {
		return castx.ToFloat64(s1, n1)
	}

	if s1, ok := r.formData[name]; ok && s1 != "" {
		return castx.ToFloat64(s1, n1)
	}

	return n1
}

func (r *Request) GetBool(name string, defaultValue ...bool) bool {
	var b1 bool

	if len(defaultValue) > 0 {
		b1 = defaultValue[0]
	}

	if b2, err := castx.ToBoolE(r.queryParams[name]); err != nil {
		return b2
	}

	if b2, err := castx.ToBoolE(r.formData[name]); err != nil {
		return b2
	}

	return b1
}

func (r *Request) GetString(name string, defaultValue ...string) string {
	var s1 string

	if len(defaultValue) > 0 {
		s1 = defaultValue[0]
	}

	if s2, ok := r.queryParams[name]; ok && s2 != "" {
		return s2
	}

	if s2, ok := r.formData[name]; ok && s2 != "" {
		return s2
	}

	return s1
}

func (r *Request) PathVariableInt(name string, defaultValue ...int) int {
	n1 := math.MinInt32

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.pathVariables[name]; ok && s1 != "" {
		return castx.ToInt(s1, n1)
	}

	return n1
}

func (r *Request) PathVariableInt64(name string, defaultValue ...int64) int64 {
	n1 := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.pathVariables[name]; ok && s1 != "" {
		return castx.ToInt64(s1, n1)
	}

	return n1
}

func (r *Request) PathVariableFloat32(name string, defaultValue ...float32) float32 {
	n1 := castx.ToFloat32(math.SmallestNonzeroFloat32)

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.pathVariables[name]; ok && s1 != "" {
		return castx.ToFloat32(s1, n1)
	}

	return n1
}

func (r *Request) PathVariableFloat64(name string, defaultValue ...float64) float64 {
	n1 := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	if s1, ok := r.pathVariables[name]; ok && s1 != "" {
		return castx.ToFloat64(s1, n1)
	}

	return n1
}

func (r *Request) PathVariableBool(name string, defaultValue ...bool) bool {
	var b1 bool

	if len(defaultValue) > 0 {
		b1 = defaultValue[0]
	}

	if b2, err := castx.ToBoolE(r.pathVariables[name]); err != nil {
		return b2
	}

	return b1
}

func (r *Request) PathVariable(name string, defaultValue ...string) string {
	var s1 string

	if len(defaultValue) > 0 {
		s1 = defaultValue[0]
	}

	if s2, ok := r.pathVariables[name]; ok && s2 != "" {
		return s2
	}

	return s1
}

func (r *Request) GetJwt(token ...string) *jwt.Token {
	var s1 string

	if len(token) > 0 {
		s1 = token[0]
	}

	if s1 == "" {
		s1 = strings.TrimSpace(r.GetHeader("Authorization"))
	}

	if s1 == "" {
		return nil
	}

	regex1 := regexp.MustCompile("[\\x20\\t]+")
	s1 = regex1.ReplaceAllString(s1, " ")

	if strings.Contains(s1, " ") {
		s1 = stringx.SubstringAfterLast(s1, " ")
	}

	if s1 == "" {
		return nil
	}

	tk, _ := securityx.ParseJsonWebToken(s1)
	return tk
}

func (r *Request) JwtClaimInt(name string, defaultValue ...int) int {
	n1 := math.MinInt32

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	token := r.GetJwt()

	if token == nil {
		return n1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return n1
	}

	if n2, err := castx.ToIntE(claims[name]); err == nil {
		return n2
	}

	return n1
}

func (r *Request) JwtClaimInt64(name string, defaultValue ...int64) int64 {
	n1 := int64(math.MinInt64)

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	token := r.GetJwt()

	if token == nil {
		return n1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return n1
	}

	if n2, err := castx.ToInt64E(claims[name]); err == nil {
		return n2
	}

	return n1
}

func (r *Request) JwtClaimFloat32(name string, defaultValue ...float32) float32 {
	n1 := castx.ToFloat32(math.SmallestNonzeroFloat32)

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	token := r.GetJwt()

	if token == nil {
		return n1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return n1
	}

	if n2, err := castx.ToFloat32E(claims[name]); err == nil {
		return n2
	}

	return n1
}

func (r *Request) JwtClaimFloat64(name string, defaultValue ...float64) float64 {
	n1 := math.SmallestNonzeroFloat64

	if len(defaultValue) > 0 {
		n1 = defaultValue[0]
	}

	token := r.GetJwt()

	if token == nil {
		return n1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return n1
	}

	if n2, err := castx.ToFloat64E(claims[name]); err == nil {
		return n2
	}

	return n1
}

func (r *Request) JwtClaimBool(name string, defaultValue ...bool) bool {
	var b1 bool

	if len(defaultValue) > 0 {
		b1 = defaultValue[0]
	}

	token := r.GetJwt()

	if token == nil {
		return b1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return b1
	}

	if b2, err := castx.ToBoolE(claims[name]); err != nil {
		return b2
	}

	return b1
}

func (r *Request) JwtClaimIntSlice(name string) []int {
	emptySlice := make([]int, 0)
	token := r.GetJwt()

	if token == nil {
		return emptySlice
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return emptySlice
	}

	if a1, err := castx.ToIntSliceE(claims[name]); err != nil {
		return a1
	}

	return emptySlice
}

func (r *Request) JwtClaimStringSlice(name string) []string {
	emptySlice := make([]string, 0)
	token := r.GetJwt()

	if token == nil {
		return emptySlice
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return emptySlice
	}

	if a1, err := castx.ToStringSliceE(claims[name]); err != nil {
		return a1
	}

	return emptySlice
}

func (r *Request) JwtClaimString(name string, defaultValue ...string) string {
	var s1 string

	if len(defaultValue) > 0 {
		s1 = defaultValue[0]
	}

	token := r.GetJwt()

	if token == nil {
		return s1
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return s1
	}

	if s2, err := castx.ToStringE(claims[name]); err == nil {
		return s2
	}

	return s1
}

func (r *Request) GetUploadedFile(formFieldName string) *multipart.FileHeader {
	return r.formFiles[formFieldName]
}

func (r *Request) GetMap(rules ...[]string) map[string]interface{} {
	var _rules []string

	if len(rules) > 0 {
		_rules = rules[0]
	}

	method := r.GetMethod()
	methods := []string{"POST", "PUT", "PATCH", "DELETE"}
	contentType := strings.ToLower(r.GetHeader("Content-Type"))
	isPostForm := strings.Contains(contentType, "application/x-www-form-urlencoded")
	isMultipartForm := strings.Contains(contentType, "multipart/form-data")
	isJson := strings.Contains(contentType, "application/json")
	isXml1 := strings.Contains(contentType, "application/xml")
	isXml2 := strings.Contains(contentType, "text/xml")
	map1 := map[string]interface{}{}

	if method == "GET" {
		for key, value := range r.queryParams {
			map1[key] = value
		}
	} else if method == "POST" && (isPostForm || isMultipartForm) {
		for key, value := range r.queryParams {
			map1[key] = value
		}

		for key, value := range r.formData {
			map1[key] = value
		}
	} else if slicex.InStringSlice(method, methods) {
		return map1
	} else if isJson {
		map1 = jsonx.MapFrom(r.GetRawBody())
	} else if isXml1 || isXml2 {
		map2 := mapx.FromXml(r.GetRawBody())

		for key, value := range map2 {
			map1[key] = value
		}
	}

	if len(map1) < 1 {
		return map[string]interface{}{}
	}

	if len(_rules) < 1 {
		return map1
	}

	ret := map[string]interface{}{}
	re1 := regexp.MustCompile(`:[^:]+$`)
	re2 := regexp.MustCompile(`:[0-9]+$`)

	for _, rule := range _rules {
		name := rule
		typ := 1
		mode := 2
		dv := ""

		if strings.HasPrefix(name, "i:") {
			name = stringx.SubstringAfter(name, ":")
			typ = 2

			if re1.MatchString(name) {
				dv = stringx.SubstringAfterLast(name, ":")
				name = re1.ReplaceAllString(name, "")
			}
		} else if strings.HasPrefix(name, "d:") {
			name = stringx.SubstringAfter(name, ":")
			typ = 3

			if re1.MatchString(name) {
				dv = stringx.SubstringAfterLast(name, ":")
				name = re1.ReplaceAllString(name, "")
			}
		} else if strings.HasPrefix(name, "s:") {
			name = stringx.SubstringAfter(name, ":")

			if re2.MatchString(name) {
				s1 := stringx.SubstringAfterLast(name, ":")
				mode = castx.ToInt(s1, 2)
				name = re2.ReplaceAllString(name, "")
			}
		} else if re2.MatchString(name) {
			s1 := stringx.SubstringAfterLast(name, ":")
			mode = castx.ToInt(s1, 2)
			name = re2.ReplaceAllString(name, "")
		}

		if strings.Contains(name, ":") {
			name = stringx.SubstringBefore(name, ":")
		}

		if name == "" {
			continue
		}

		switch typ {
		case 1:
			value := castx.ToString(map1[name])

			switch mode {
			case 1, 2:
				value = stringx.StripTags(value)
			}

			ret[name] = value
		case 2:
			var value int

			if n1, err := castx.ToIntE(dv); err == nil {
				value = castx.ToInt(map1[name], n1)
			} else {
				value = castx.ToInt(map1[name])
			}

			ret[name] = value
		case 3:
			var value float64

			if n1, err := castx.ToFloat64E(dv); err == nil {
				value = castx.ToFloat64(map1[name], n1)
			} else {
				value = castx.ToFloat64(map1[name])
			}

			ret[name] = numberx.ToDecimalString(value)
		}
	}

	return ret
}

func (r *Request) GetRawBody() []byte {
	return r.rawBody
}

func (r *Request) GetRouteRule() *mvc.RouteRule {
	return r.routeRule
}

func (r *Request) WithMiddleware(m Middleware) *Request {
	r.middlewares = append(r.middlewares, m)
	return r
}

func (r *Request) WithMiddlewares(entries []Middleware) *Request {
	if len(entries) > 0 {
		r.middlewares = append(r.middlewares, entries...)
	}

	return r
}

func (r *Request) GetMiddlewares() []Middleware {
	return r.middlewares
}

func (r *Request) GetExecStart() time.Time {
	return r.execStart
}

func (r *Request) Next(flag ...bool) bool {
	if len(flag) > 0 {
		r.next = flag[0]
		return false
	}

	return r.next
}
