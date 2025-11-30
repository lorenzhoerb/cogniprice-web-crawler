package service

import (
	"errors"
	"log"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
	"github.com/lorenzhoerb/cogniprice/shared"
)

// JobRepository defines methods to manage jobs in the scheduler service.
type JobRepository interface {
	// ListDuoJobs returns up to 'limit' duo jobs.
	// If limit == 0, all duo jobs are returned.
	//GetDueJobs(limit int) ([]*model.Job, error)

	// GetJob retrieves a job by its ID.
	//GetJob(id int) (*model.Job, error)
	GetByID(id int) (*model.Job, error)

	// GetByURL retrieves a job by its URL.
	// Returns null if not found.
	GetByURL(url string) (*model.Job, error)

	// List all jobs and filters them
	List(filter *model.ListJobsFilter) ([]*model.Job, *shared.Pagination, error)

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

func (s *JobService) GetJob(id int) (*model.JobResponse, error) {
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
