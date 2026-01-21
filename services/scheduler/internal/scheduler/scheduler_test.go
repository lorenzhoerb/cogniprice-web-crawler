package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/config"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/mocks"
	"go.uber.org/mock/gomock"
)

func TestSchedulerDispatchSuccessfully(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dispatcher := mocks.NewMockDispatcher(ctrl)
	repo := mocks.NewMockSchedulerJobRepository(ctrl)
	batchSize := 5

	scheduler := NewScheduler(
		&config.SchedulerConfig{Interval: 50 * time.Microsecond, BatchSize: batchSize},
		repo,
		dispatcher,
	)

	jobs := []*model.Job{
		{ID: 1},
		{ID: 2},
	}

	// set expectations BEFORE calling Run
	repo.EXPECT().GetDue(batchSize).Return(jobs, nil).Times(1)
	repo.EXPECT().SaveAll(gomock.Any()).Do(func(jobs []*model.Job) {
		if len(jobs) != 2 {
			t.Fatalf("expected 2 jobs, got %d", len(jobs))
		}
		if jobs[0].Status != model.JobStatusInProgress {
			t.Fatalf("job 0 status not updated")
		}
		if jobs[0].DispatchedAt == nil {
			t.Fatalf("job 0 dispatchedAt not set")
		}
	}).Return(nil).Times(1)

	dispatcher.EXPECT().DispatchJobs(gomock.Any()).Return(nil).Times(1)

	// After first batch, GetDue returns empty slice
	repo.EXPECT().GetDue(batchSize).Return(nil, nil).Times(1)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Microsecond)
	defer cancel()

	// run scheduler (will stop when context is done)
	scheduler.Run(ctx)
}

func TestSchedulerIdleGetDuo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dispatcher := mocks.NewMockDispatcher(ctrl)
	repo := mocks.NewMockSchedulerJobRepository(ctrl)
	batchSize := 5

	scheduler := NewScheduler(
		&config.SchedulerConfig{Interval: 1 * time.Millisecond, BatchSize: batchSize},
		repo,
		dispatcher,
	)

	repo.EXPECT().GetDue(batchSize).Return(nil, nil).MaxTimes(510).MinTimes(490)
	repo.EXPECT().SaveAll(gomock.Any()).Times(0)
	dispatcher.EXPECT().DispatchJobs(gomock.Any()).Times(0)

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// run scheduler (will stop when context is done)
	scheduler.Run(ctx)
}
