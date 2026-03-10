package promptrails

import (
	"errors"
	"testing"
)

func TestNewErrorForStatus(t *testing.T) {
	tests := []struct {
		status int
		check  func(error) bool
	}{
		{400, func(e error) bool { var v *ValidationError; return errors.As(e, &v) }},
		{401, func(e error) bool { var v *UnauthorizedError; return errors.As(e, &v) }},
		{402, func(e error) bool { var v *QuotaExceededError; return errors.As(e, &v) }},
		{403, func(e error) bool { var v *ForbiddenError; return errors.As(e, &v) }},
		{404, func(e error) bool { var v *NotFoundError; return errors.As(e, &v) }},
		{429, func(e error) bool { var v *RateLimitError; return errors.As(e, &v) }},
		{500, func(e error) bool { var v *ServerError; return errors.As(e, &v) }},
		{503, func(e error) bool { var v *ServerError; return errors.As(e, &v) }},
	}

	for _, tt := range tests {
		err := newErrorForStatus(tt.status, "test error", "TEST", nil)
		if !tt.check(err) {
			t.Errorf("status %d: unexpected error type %T", tt.status, err)
		}
	}
}

func TestAPIErrorMessage(t *testing.T) {
	err := newErrorForStatus(404, "agent not found", "NOT_FOUND", nil)
	if err.Error() != "agent not found (code=NOT_FOUND, status=404)" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
