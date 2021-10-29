package httpx

type HttpError struct {
	statusCode int
}

func NewHttpError(statusCode int) HttpError {
	if statusCode < 400 {
		statusCode = 500
	}

	return HttpError{statusCode: statusCode}
}

func (p HttpError) GetContentType() string {
	return ""
}

func (p HttpError) GetContents() interface{} {
	return ""
}

func (p HttpError) GetStatusCode() int {
	return p.statusCode
}
