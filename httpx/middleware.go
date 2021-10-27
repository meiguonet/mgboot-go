package httpx

type Middleware interface {
	GetType() int
	GetOrder() int
	PreHandle(request Request, response Response)
	PostHandle(request Request, response Response)
}
