package httpx

type Middleware interface {
	GetName() string
	GetType() int
	GetOrder() int
	PreHandle(request *Request, response *Response)
	PostHandle(request *Request, response *Response)
}
