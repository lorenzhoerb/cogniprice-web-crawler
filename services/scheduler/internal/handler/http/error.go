package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIError struct {
	Message string       `json:"message"`
	Code    string       `json:"code"`
	Errors  []FieldError `json:"errors,omitempty"`
}

var ErrInternalServer = APIError{
	Message: "internal server error",
	Code:    "INTERNAL_SERVER_ERROR",
}

// ErrorHandler is a middleware that handles service-level errors and validation errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// If no errors were collected, nothing to do
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			handleValidationErrors(c, verrs)
			return
		}

		var appErr *service.AppError
		if errors.As(err, &appErr) {
			handleAppErrors(c, appErr)
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, &ErrInternalServer)
	}
}

func handleValidationErrors(c *gin.Context, verrs validator.ValidationErrors) {
	apiErr := &APIError{
		Message: "one or more fields are invalid",
		Code:    "INVALID_REQUEST",
		Errors:  []FieldError{},
	}

	for _, fe := range verrs {
		apiErr.Errors = append(apiErr.Errors, FieldError{
			Field:   fe.Field(), // or JSON name if you register tag name func
			Message: validationErrorMsg(fe),
		})
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, apiErr)
}

func handleAppErrors(c *gin.Context, appErr *service.AppError) {
	var fieldErrors []FieldError
	for _, fe := range appErr.Errors {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   fe.Field,
			Message: fe.Message,
		})
	}
	c.AbortWithStatusJSON(appErr.Status, APIError{
		Message: appErr.Message,
		Code:    appErr.Code,
		Errors:  fieldErrors,
	})
}

// validationErrorMsg returns a user-friendly message for a FieldError
func validationErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "url":
		return "must be a valid URL"
	case "email":
		return "must be a valid email"
	case "min":
		return "must be at least " + fe.Param() + " characters long"
	case "max":
		return "must be at most " + fe.Param() + " characters long"
	case "interval":
		return "interval must be a valid duration (e.g., '10s', '5m', '1h') and at least 1 hour"
	case "jobstatus":
		return "status must be one of: scheduled, in_progress, paused, failed"
	default:
		return fe.Error()
	}
}
