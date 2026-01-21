package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/shared/logger"
	"go.uber.org/zap"
)

// Dispatcher handles job submissions to the worker queue.

//go:generate mockgen -destination=../../mocks/mock_dispatcher.go -package=mocks . Dispatcher
type Dispatcher interface {
	// Dispatches all jobs as a batch to the worker queue.
	DispatchJobs(jobs []model.JobDispatched) error
}

//go:generate mockgen -destination=../../mocks/mock_scheduler_job_repository.go -mock_names=JobRepository=MockSchedulerJobRepository -package=mocks . JobRepository
type JobRepository interface {
	// ListDue returns up to 'limit' duo jobs.
	// If limit == 0, all duo jobs are returned.
	GetDue(limit int) ([]*model.Job, error)

	// UpdateJobs batch updates all jobs specified
	SaveAll(job []*model.Job) error
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
	batchSize := cfg.BatchSize
	if batchSize <= 0 {
		batchSize = 100
	}

	return &Scheduler{
		Repo:       repo,
		Interval:   cfg.Interval,
		BatchSize:  batchSize,
		Dispatcher: dispatcher,
	}
}

// Start starts the job cycle
func (s *Scheduler) Run(ctx context.Context) {
	logger.Log.Info("scheduler started",
		zap.String("timeInterval", s.Interval.String()),
		zap.Int("batchSize", s.BatchSize))

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("scheduler stopped due to context cancellation")
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
	logger.Log.Debug("checking for due jobs")
	dueJobs, err := s.Repo.GetDue(s.BatchSize)
	if err != nil {
		logger.Log.Error("failed to fetch due jobs", zap.Error(err))
		return fmt.Errorf("get due jobs failed: %w", err)
	}

	if len(dueJobs) == 0 {

		// no due jobs to dispatch
		return nil
	}

	logger.Log.Info("found due jobs to dispatch",
		zap.Int("jobCount", len(dueJobs)),
	)

	dispatchedAt := time.Now()

	var jobsDispatched []model.JobDispatched
	for _, job := range dueJobs {
		// set job metadata
		job.Status = model.JobStatusInProgress
		job.DispatchedAt = &dispatchedAt

		jobsDispatched = append(jobsDispatched, model.JobDispatched{
			ID:           job.ID,
			URL:          job.URL,
			DispatchedAt: dispatchedAt,
		})
	}

	// update job metadata
	if err := s.Repo.SaveAll(dueJobs); err != nil {
		return fmt.Errorf("failed to update job status to DISPATCHED: %w", err)
	}

	// dispatch jobs to worker queue
	if err := s.Dispatcher.DispatchJobs(jobsDispatched); err != nil {
		// TODO: Rollback
		return fmt.Errorf("failed to dispatch jobs: %w", err)
	}

	logger.Log.Info("successfully dispatched jobs",
		zap.Int("jobCount", len(jobsDispatched)),
	)

	return nil
}

func (s *Scheduler) handleDispatchFail() {
}
