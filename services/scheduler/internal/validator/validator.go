package validator

import (
	"time"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
)

// interval must be at least 1 hour
var interval validator.Func = func(fl validator.FieldLevel) bool {
	intervalStr, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	d, err := time.ParseDuration(intervalStr)
	if err != nil {
		return false
	}
	if d < time.Hour {
		return false
	}
	return true
}

// interval must be at least 1 hour
var jobStatus validator.Func = func(fl validator.FieldLevel) bool {
	status, ok := fl.Field().Interface().(model.JobStatus)
	if !ok {
		return false
	}
	return status.IsValid()
}

// interval must be at least 1 hour
var sortOrder validator.Func = func(fl validator.FieldLevel) bool {
	order, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return order == "asc" || order == "desc"
}

var jobSortCol validator.Func = func(fl validator.FieldLevel) bool {
	col, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return col == "created_at" || col == "url" || col == "next_run_at" || col == "status"
}

func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("interval", interval)
		v.RegisterValidation("jobstatus", jobStatus)
		v.RegisterValidation("sortorder", sortOrder)
		v.RegisterValidation("jobsortcol", jobSortCol)
	}
}
