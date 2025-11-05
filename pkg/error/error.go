package error

type HTTPError struct {
	Code    int
	Message string
}

func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Message: message}
}
