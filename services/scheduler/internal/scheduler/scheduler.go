package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
)

// Dispatcher handles job submissions to the worker queue.
//
//go:generate mockgen -destination=../../mocks/mock_dispatcher.go -package=mocks github.com/lorenzhoerb/cogniprice/services/scheduler/internal/scheduler Dispatcher
type Dispatcher interface {
	// Dispatches all jobs as a batch to the worker queue.
	DispatchJobs(jobs []model.JobDispatched) error
}

//go:generate mockgen -destination=../../mocks/scheduler_job_repository.go -package=mocks github.com/lorenzhoerb/cogniprice/services/scheduler/internal/scheduler JobRepository
type JobRepository interface {
	// ListDuoJobs returns up to 'limit' duo jobs.
	// If limit == 0, all duo jobs are returned.
	GetDueJobs(limit int) ([]*model.Job, error)

	// UpdateJobs batch updates all jobs specified
	UpdateJobs(job []*model.Job) error
}

// Scheduler manages the periodic dispatching of due jobs to the worker queue.
type Scheduler struct {
	// Repo provides access to job storage for retrieving and updating job states.
	Repo JobRepository

	// Dispatcher handles the submission of jobs to the worker queue.
	Dispatcher Dispatcher

	// Interval defines how often the scheduler checks for due jobs.
	Interval time.Duration

	// BatchSize specifies the maximum number of jobs to schedule in a single run.
	BatchSize int
}

func NewScheduler(cfg *config.SchedulerConfig, repo JobRepository, dispatcher Dispatcher) *Scheduler {
	return &Scheduler{
		Repo:       repo,
		Interval:   cfg.Interval,
		Dispatcher: dispatcher,
	}
}

// Start starts the job cycle
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.dispatchDueJobs(); err != nil {
				s.handleDispatchFail()
			}
		}
	}
}

// dispatchDueJobs dispatches jobs due.
// Upon dispatching it ensures that the job status is set to dispatched.
func (s *Scheduler) dispatchDueJobs() error {
	log.Println("dispatching due jobs")
	dueJobs, err := s.Repo.GetDueJobs(s.BatchSize)
	if err != nil {
		return fmt.Errorf("get due jobs failed: %w", err)
	}

	if len(dueJobs) == 0 {
		// no due jobs to dispatch
		return nil
	}

	dispatchedAt := time.Now()

	var jobsDispatched []model.JobDispatched
	for _, job := range dueJobs {
		// set job metadata
		job.Status = model.JobStatus(model.InProgress)
		job.DispatchedAt = &dispatchedAt

		jobsDispatched = append(jobsDispatched, model.JobDispatched{
			ID:           job.ID,
			URL:          job.URL,
			DispatchedAt: dispatchedAt,
		})
	}

	// update job metadata
	if err := s.Repo.UpdateJobs(dueJobs); err != nil {
		return fmt.Errorf("failed to update job status to DISPATCHED: %w", err)
	}

	// dispatch jobs to worker queue
	if err := s.Dispatcher.DispatchJobs(jobsDispatched); err != nil {
		// TODO: Rollback
		return fmt.Errorf("failed to dispatch jobs: %w", err)
	}

	return nil
}

func (s *Scheduler) handleDispatchFail() {
}
