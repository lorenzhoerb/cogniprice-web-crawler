package service

import "fmt"

type FileError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type AppError struct {
	Message  string      `json:"message"`
	Code     string      `json:"code"`
	Status   int         `json:"status"`
	SvcError error       `json:"-"`
	Errors   []FileError `json:"errors,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrJobWithURLExists = &AppError{
		Message: "a job with the given URL already exists",
		Code:    "JOB_WITH_URL_EXISTS",
		Status:  409,
	}
	ErrCannotPauseJob = &AppError{
		Message: "cannot pause job in current state",
		Code:    "CANNOT_PAUSE_JOB",
		Status:  400,
	}
)

func ErrNotFound(id any) *AppError {
	return &AppError{
		Message: fmt.Sprintf("job with id %v not found", id),
		Code:    "NOT_FOUND",
		Status:  404,
	}
}

func ErrInvalidField(field, msg string) *AppError {
	return &AppError{
		Message: fmt.Sprintf("invalid value for field '%s'", field),
		Code:    "INVALID_FIELD",
		Status:  400,
		Errors: []FileError{
			{
				Field:   field,
				Message: msg,
			},
		},
	}
}
