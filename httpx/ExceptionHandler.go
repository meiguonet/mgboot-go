package httpx

type ExceptionHandler interface {
	GetExceptionName() string
	MatchException(err error) bool
	HandleException(err error) ResponsePayload
}
