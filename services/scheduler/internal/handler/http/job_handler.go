package http

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service"
)

type JobHandler struct {
	Svc service.JobService
}

func NewJobHandler(svc *service.JobService) *JobHandler {
	return &JobHandler{
		Svc: *svc,
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

	jobResp, err := h.Svc.GetJob(id)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, jobResp)
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
