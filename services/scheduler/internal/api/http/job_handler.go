package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
)

//go:generate mockgen -destination=../../../mocks/mock_job_service.go -package=mocks github.com/lorenzhoerb/cogniprice/services/scheduler/internal/api/http JobService
type JobService interface {
	CreateJob(*model.CreateJobRequest) (*model.JobResponse, error)
	GetJobByID(int) (*model.JobResponse, error)
	ListJobs(*model.ListJobsFilter) (*model.PaginatedJobsResponse, error)
	PauseJob(int) (*model.JobResponse, error)
	ResumeJob(int) (*model.JobResponse, error)
	DeleteJob(int) error
}

type JobHandler struct {
	Svc JobService
}

func NewJobHandler(svc JobService) *JobHandler {
	return &JobHandler{
		Svc: svc,
	}
}

// CreateJob validates the request (including ensuring interval >= 1 hour)
// and returns a created job placeholder. In a full implementation this
// would call the service layer to persist the job.
func (h *JobHandler) CreateJob(c *gin.Context) {
	var req model.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	jobResp, err := h.Svc.CreateJob(&req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, jobResp)
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, err := parseJobID(c)
	if err != nil {
		c.Error(err)
		return
	}

	jobResp, err := h.Svc.GetJobByID(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, jobResp)
}

// ListJobs retrieves jobs and returns it as a pagination.
// Jobs can be filtered by URL substring, and Status
func (h *JobHandler) ListJobs(c *gin.Context) {
	var filter model.ListJobsFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.Error(err)
		return
	}

	paginatedJobs, err := h.Svc.ListJobs(&filter)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, paginatedJobs)
}

func (h *JobHandler) PauseJob(c *gin.Context) {
	id, err := parseJobID(c)
	if err != nil {
		c.Error(err)
		return
	}

	jobResp, err := h.Svc.PauseJob(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, jobResp)
}

func (h *JobHandler) ResumeJob(c *gin.Context) {
	id, err := parseJobID(c)
	if err != nil {
		c.Error(err)
		return
	}

	jobResp, err := h.Svc.ResumeJob(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, jobResp)
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id, err := parseJobID(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.Svc.DeleteJob(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}

func parseJobID(c *gin.Context) (int, error) {
	id := c.Param("id")
	idNum, err := strconv.Atoi(id)
	if err != nil {
		return 0, &service.AppError{
			Message: fmt.Sprintf("invalid job id: %s", id),
			Code:    "INVALID_JOB_ID",
			Status:  400,
		}
	}
	return idNum, nil
}
