package error

// HTTPError represents an application-level HTTP error
type HTTPError struct {
	Code    int
	Message string
	Error   error
}

// NewHTTPError creates a new HTTPError instance
func NewHTTPError(code int, message string, error error) *HTTPError {
	return &HTTPError{Code: code, Message: message, Error: error}
}

// GetFullMessage returns the combined message and error text
func (h *HTTPError) GetFullMessage() string {
	return h.Message + ": " + h.Error.Error()
}
