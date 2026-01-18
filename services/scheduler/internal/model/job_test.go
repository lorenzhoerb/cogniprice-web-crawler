package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduleNextRun(t *testing.T) {
	tolerance := 10 * time.Millisecond // allow a small timing difference

	now := time.Now()

	tests := []struct {
		name           string
		lastCrawledAt  *time.Time
		interval       time.Duration
		expectedStatus JobStatus
		expectedNext   func() time.Time // function to compute expected NextRunAt
	}{
		{
			name:           "never crawled, schedules now",
			lastCrawledAt:  nil,
			interval:       1 * time.Hour,
			expectedStatus: JobStatusScheduled,
			expectedNext: func() time.Time {
				return time.Now()
			},
		},
		{
			name:           "scheduled in the future",
			lastCrawledAt:  ptrTime(now),
			interval:       1 * time.Hour,
			expectedStatus: JobStatusScheduled,
			expectedNext: func() time.Time {
				return now.Add(1 * time.Hour)
			},
		},
		{
			name:           "overdue job, next run should be now",
			lastCrawledAt:  ptrTime(now.Add(-2 * time.Hour)),
			interval:       1 * time.Hour,
			expectedStatus: JobStatusScheduled,
			expectedNext: func() time.Time {
				return time.Now()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				LastCrawledAt: tt.lastCrawledAt,
				Interval:      tt.interval,
			}

			job.ScheduleNextRun()

			// Check status
			assert.Equal(t, tt.expectedStatus, job.Status, "job status mismatch")

			// Check NextRunAt within tolerance
			diff := job.NextRunAt.Sub(tt.expectedNext())
			if diff < 0 {
				diff = -diff
			}
			assert.True(t, diff <= tolerance, "NextRunAt differs by %v, exceeds tolerance %v", diff, tolerance)
		})
	}
}

func TestJobPause(t *testing.T) {
	tests := []struct {
		name         string
		initialJob   Job
		wantStatus   JobStatus
		wantPauseReq bool
		wantError    error
	}{
		{
			name: "pause job in scheduled state",
			initialJob: Job{
				Status:         JobStatusScheduled,
				PauseRequested: false,
			},
			wantStatus:   JobStatusPaused,
			wantPauseReq: false,
			wantError:    nil,
		},
		{
			name: "pause job in in_progress state",
			initialJob: Job{
				Status:         JobStatusInProgress,
				PauseRequested: false,
			},
			wantStatus:   JobStatusInProgress,
			wantPauseReq: true,
			wantError:    nil,
		},
		{
			name: "pause job in paused state",
			initialJob: Job{
				Status:         JobStatusPaused,
				PauseRequested: false,
			},
			wantStatus:   JobStatusPaused,
			wantPauseReq: false,
			wantError:    nil,
		},
		{
			name: "pause job in failed state",
			initialJob: Job{
				Status:         JobStatusFailed,
				PauseRequested: false,
			},
			wantStatus:   JobStatusFailed,
			wantPauseReq: false,
			wantError:    ErrCannotPause,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := tt.initialJob

			err := job.Pause()

			assert.Equal(t, tt.wantError, err, "error mismatch")
			assert.Equal(t, tt.wantStatus, job.Status, "status mismatch")
			assert.Equal(t, tt.wantPauseReq, job.PauseRequested, "PauseRequested mismatch")

		})
	}

}

func TestResumeJob(t *testing.T) {
	tolerance := 10 * time.Millisecond
	now := time.Now()

	tests := []struct {
		name           string
		initialStatus  JobStatus
		lastCrawledAt  *time.Time
		interval       time.Duration
		expectedStatus JobStatus
		expectRetry    int
		expectNextRun  func() time.Time // computes expected NextRunAt
	}{
		{
			name:           "paused job resumes",
			initialStatus:  JobStatusPaused,
			lastCrawledAt:  nil,
			interval:       1 * time.Hour,
			expectedStatus: JobStatusScheduled,
			expectRetry:    0,
			expectNextRun:  func() time.Time { return time.Now() },
		},
		{
			name:           "failed job resumes",
			initialStatus:  JobStatusFailed,
			lastCrawledAt:  ptrTime(now.Add(-1 * time.Hour)),
			interval:       30 * time.Minute,
			expectedStatus: JobStatusScheduled,
			expectRetry:    0,
			expectNextRun:  func() time.Time { return time.Now() },
		},
		{
			name:           "already scheduled job does nothing",
			initialStatus:  JobStatusScheduled,
			lastCrawledAt:  ptrTime(now),
			interval:       1 * time.Hour,
			expectedStatus: JobStatusScheduled,
			expectRetry:    5,                               // assume job had retries
			expectNextRun:  func() time.Time { return now }, // should not change
		},
		{
			name:           "in progress job does nothing",
			initialStatus:  JobStatusInProgress,
			lastCrawledAt:  ptrTime(now),
			interval:       1 * time.Hour,
			expectedStatus: JobStatusInProgress,
			expectRetry:    2,                               // assume job had retries
			expectNextRun:  func() time.Time { return now }, // should not change
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := &Job{
				Status:        tt.initialStatus,
				LastCrawledAt: tt.lastCrawledAt,
				Interval:      tt.interval,
				RetryAttempts: tt.expectRetry,
			}

			err := job.Resume()
			assert.NoError(t, err)

			// Status check
			assert.Equal(t, tt.expectedStatus, job.Status, "status mismatch")

			// RetryAttempts check
			if tt.initialStatus == JobStatusInProgress || tt.initialStatus == JobStatusScheduled {
				assert.Equal(t, tt.expectRetry, job.RetryAttempts, "retry attempts should remain unchanged")
			} else {
				assert.Equal(t, 0, job.RetryAttempts, "retry attempts should be reset")

			}

			// NextRunAt check (only for resumed jobs)
			if job.Status == JobStatusScheduled {
				diff := job.NextRunAt.Sub(tt.expectNextRun())
				if diff < 0 {
					diff = -diff
				}
				assert.True(t, diff <= tolerance, "NextRunAt differs by %v, exceeds tolerance %v", diff, tolerance)
			}
		})
	}
}

// helper to get pointer to a time
func ptrTime(t time.Time) *time.Time {
	return &t
}
