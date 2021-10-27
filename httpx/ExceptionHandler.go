package httpx

type ExceptionHandler interface {
	MatchException(err error) bool
	HandleException(err error) ResponsePayload
}
