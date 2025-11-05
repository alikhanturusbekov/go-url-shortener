package error

type HTTPError struct {
	Code    int
	Message string
	Error   error
}

func NewHTTPError(code int, message string, error error) *HTTPError {
	return &HTTPError{Code: code, Message: message, Error: error}
}

func (h *HTTPError) GetFullMessage() string {
	return h.Message + ": " + h.Error.Error()
}
