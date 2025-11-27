package model

import "time"

type CreateJobRequest struct {
	URL      string `json:"url" binding:"required,url"`
	Interval string `json:"interval" binding:"required,interval"`
}

type JobResponse struct {
	ID             uint       `json:"id"`
	URL            string     `json:"url"`
	Status         string     `json:"status"`
	Interval       string     `json:"interval"`
	RetryAttempts  int        `json:"retryAttempts"`
	PauseRequested bool       `json:"pauseRequested"`
	DispatchedAt   *time.Time `json:"dispatchedAt"`
	NextRunAt      *time.Time `json:"nextRunAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type JobFilter struct {
	URL    *string
	status *string
}

type PaginatedJobsResponse struct {
	Jobs       []JobResponse `json:"jobs"`
	TotalCount int64         `json:"totalCount"`
	TotalPages int           `json:"totalPages"`
	Page       int           `json:"page"`
	PageSize   int           `json:"pageSize"`
}

func ToJobResponse(j *Job) *JobResponse {
	return &JobResponse{
		ID:             j.ID,
		URL:            j.URL,
		Status:         j.Status.String(),
		Interval:       j.Interval.String(),
		RetryAttempts:  j.RetryAttempts,
		PauseRequested: j.PauseRequested,
		DispatchedAt:   j.DispatchedAt,
		NextRunAt:      &j.NextRunAt,
		CreatedAt:      j.CreatedAt,
		UpdatedAt:      j.UpdatedAt,
	}
}
