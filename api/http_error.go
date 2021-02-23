package api

type HTTPError struct {
	Code    int
	Message string
}

func NewHTTPError(code int, msg string) HTTPError {
	return HTTPError{
		Code:    code,
		Message: msg,
	}
}

func (e HTTPError) Error() string {
	return e.Message
}
