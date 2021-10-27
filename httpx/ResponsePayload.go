package httpx

type ResponsePayload interface {
	GetContentType() string
	GetContents() interface{}
}
