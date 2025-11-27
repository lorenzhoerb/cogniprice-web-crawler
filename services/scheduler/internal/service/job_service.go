package service

import (
	"log"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
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

	// PutJob inserts or updates a job.
	// If job.ID is empty, an ID is generated and assigned to the same object.
	Save(job *model.Job) error
}

type JobService struct {
	Repo JobRepository
}

// NewJobService instantiates a JobService
func NewJobService(repo JobRepository) *JobService {
	return &JobService{
		Repo: repo,
	}
}

func (s *JobService) CreateJob(req *model.CreateJobRequest) (*model.JobResponse, error) {
	log.Printf("Creating job with URL: %s and Interval: %s\n", req.URL, req.Interval)
	interval, _ := time.ParseDuration(req.Interval) // already validated

	// Check for existing job with the same URL
	existingJob, err := s.Repo.GetByURL(req.URL)
	if err != nil {
		return nil, err
	}

	if existingJob != nil {
		return nil, ErrJobWithURLExists
	}

	job := &model.Job{
		URL:       req.URL,
		Interval:  interval,
		Status:    model.Scheduled,
		NextRunAt: time.Now(),
	}

	err = s.Repo.Save(job)
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

	err = s.Repo.Save(job)
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

	err = s.Repo.Save(job)
	if err != nil {
		return nil, err
	}

	return model.ToJobResponse(job), nil
}

func (s *JobService) getJobByIDOrNotFound(id int) (*model.Job, error) {
	job, err := s.Repo.GetByID(id)
	if err == nil {
		return job, nil
	}

	if err == repository.ErrNotFound {
		return nil, ErrNotFound(id)
	}

	return nil, err
}
