package promptrails

import "fmt"

// APIError is the base error returned by the API.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
	Details    map[string]any
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
	}
	return fmt.Sprintf("%s (status=%d)", e.Message, e.StatusCode)
}

// Typed errors for specific HTTP status codes.

// ValidationError is returned on 400 Bad Request.
type ValidationError struct{ APIError }

// UnauthorizedError is returned on 401 Unauthorized.
type UnauthorizedError struct{ APIError }

// ForbiddenError is returned on 403 Forbidden.
type ForbiddenError struct{ APIError }

// NotFoundError is returned on 404 Not Found.
type NotFoundError struct{ APIError }

// QuotaExceededError is returned on 402 Payment Required.
type QuotaExceededError struct{ APIError }

// RateLimitError is returned on 429 Too Many Requests.
type RateLimitError struct{ APIError }

// ServerError is returned on 5xx errors.
type ServerError struct{ APIError }

func newErrorForStatus(statusCode int, msg, code string, details map[string]any) error {
	base := APIError{StatusCode: statusCode, Message: msg, Code: code, Details: details}
	switch {
	case statusCode == 400:
		return &ValidationError{base}
	case statusCode == 401:
		return &UnauthorizedError{base}
	case statusCode == 402:
		return &QuotaExceededError{base}
	case statusCode == 403:
		return &ForbiddenError{base}
	case statusCode == 404:
		return &NotFoundError{base}
	case statusCode == 429:
		return &RateLimitError{base}
	case statusCode >= 500:
		return &ServerError{base}
	default:
		return &base
	}
}
