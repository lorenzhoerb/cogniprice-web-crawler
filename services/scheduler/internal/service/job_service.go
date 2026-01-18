package service

import (
	"errors"
	"log"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
	"github.com/lorenzhoerb/cogniprice/shared/pagination"
)

// JobRepository defines methods to manage jobs in the scheduler service.

//go:generate mockgen -destination=../../mocks/mock_job_repository_service.go -package=mocks github.com/lorenzhoerb/cogniprice/services/scheduler/internal/service JobRepository
type JobRepository interface {

	// GetJobByID retrieves a job by its ID. Returns null and NotFound error if not found.
	GetByID(id int) (*model.Job, error)

	// GetByURL retrieves a job by its URL.
	// Returns null if not found.
	GetByURL(url string) (*model.Job, error)

	// List all jobs and filters them
	List(filter *model.ListJobsFilter) ([]*model.Job, *pagination.Pagination, error)

	// PutJob inserts or updates a job.
	// If job.ID is empty, an ID is generated and assigned to the same object.
	Save(job *model.Job) error

	// Delete removes a job by its ID.
	Delete(id int) error
}

type JobService struct {
	repo JobRepository
}

// NewJobService instantiates a JobService
func NewJobService(repo JobRepository) *JobService {
	return &JobService{
		repo: repo,
	}
}

func (s *JobService) CreateJob(req *model.CreateJobRequest) (*model.JobResponse, error) {
	log.Printf("Creating job with URL: %s and Interval: %s\n", req.URL, req.Interval)
	interval, _ := time.ParseDuration(req.Interval) // already validated

	// Check for existing job with the same URL
	_, err := s.repo.GetByURL(req.URL)
	if err == nil {
		// Job exists → cannot create duplicate
		return nil, ErrJobWithURLExists
	}
	if !errors.Is(err, repository.ErrNotFound) {
		// Some other error occurred
		return nil, err
	}

	job := &model.Job{
		URL:       req.URL,
		Interval:  interval,
		Status:    model.JobStatusScheduled,
		NextRunAt: time.Now(),
	}

	err = s.repo.Save(job)
	if err != nil {
		return nil, err
	}

	return model.ToJobResponse(job), nil
}

func (s *JobService) GetJobByID(id int) (*model.JobResponse, error) {
	log.Printf("Retrieving job with ID: %d\n", id)
	job, err := s.getJobByIDOrNotFound(id)
	if err != nil {
		return nil, err
	}

	return model.ToJobResponse(job), nil
}

func (s *JobService) ListJobs(filter *model.ListJobsFilter) (*model.PaginatedJobsResponse, error) {
	log.Printf("List jobs %+v\n", filter)

	jobs, pagination, err := s.repo.List(filter)
	if err != nil {
		return nil, err
	}

	// Convert to JobResponse
	jobResponse := make([]*model.JobResponse, 0, len(jobs))
	for _, job := range jobs {
		jobResponse = append(jobResponse, model.ToJobResponse(job))
	}

	return &model.PaginatedJobsResponse{
		Items:      jobResponse,
		TotalCount: pagination.Total,
		TotalPages: pagination.TotalPages(),
		Page:       pagination.CurrentPage(),
		PageSize:   pagination.PageSize,
	}, nil
}

// PauseJob attempts to pause the job with the given ID.
//
// Behavior:
//   - If the job does not exist, it returns ErrNotFound.
//   - If the job is in Scheduled state, it is immediately paused
//     and the status becomes Paused.
//   - If the job is InProgress, the job continues running, the
//     PauseRequested flag is set to true, and the status remains InProgress.
//   - If the job is in a Failed state, pausing is not allowed and
//     ErrCannotPauseJob is returned.
//
// Side effects:
//   - Successfully paused or pause-requested jobs are persisted.
//   - Jobs that cannot be paused are not saved.
//
// Returns:
//   - A JobResponse reflecting the current job state.
//   - An error if the job cannot be paused or does not exist.
func (s *JobService) PauseJob(id int) (*model.JobResponse, error) {
	log.Printf("Pausing job with ID: %d\n", id)
	job, err := s.getJobByIDOrNotFound(id)
	if err != nil {
		return nil, err
	}

	err = job.Pause()
	if err != nil {
		return nil, ErrCannotPauseJob
	}

	err = s.repo.Save(job)
	if err != nil {
		return nil, err
	}

	return model.ToJobResponse(job), nil
}

func (s *JobService) ResumeJob(id int) (*model.JobResponse, error) {
	log.Printf("Resuming job with ID: %d\n", id)
	job, err := s.getJobByIDOrNotFound(id)
	if err != nil {
		return nil, err
	}

	job.Resume()

	err = s.repo.Save(job)
	if err != nil {
		return nil, err
	}

	return model.ToJobResponse(job), nil
}

func (s *JobService) getJobByIDOrNotFound(id int) (*model.Job, error) {
	job, err := s.repo.GetByID(id)
	if err == nil {
		return job, nil
	}

	if err == repository.ErrNotFound {
		return nil, ErrNotFound(id)
	}

	return nil, err
}

func (s *JobService) DeleteJob(id int) error {
	log.Printf("Deleting job with ID: %d\n", id)
	_, err := s.getJobByIDOrNotFound(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
