package BuiltinException

import "fmt"

type HttpErrorException struct {
	statusCode int
}

func NewHttpErrorException(statusCode int) HttpErrorException {
	if statusCode < 400 {
		statusCode = 500
	}

	return HttpErrorException{statusCode: statusCode}
}

func (ex HttpErrorException) Error() string {
	return fmt.Sprintf("http error %d", ex.statusCode)
}

func (ex HttpErrorException) GetStatusCode() int {
	return ex.statusCode
}
