package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var ErrCannotPause = errors.New("cannot pause job in current state")

// JobStatus represents the current status of a job.
//
// Status values:
//   - Scheduled: The job is ready to be dispatched if NextRunAt <= now.
//   - Running: The job was dispatched successfully and is currently executing.
//   - Paused: The job will not be scheduled until resumed.
//   - Failed: The job exceeded its failure/retry limit and is no longer eligible for dispatch.
type JobStatus string

const (
	JobStatusScheduled  JobStatus = "scheduled"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusPaused     JobStatus = "paused"
	JobStatusFailed     JobStatus = "failed"
)

func (s JobStatus) IsValid() bool {
	switch s {
	case JobStatusScheduled, JobStatusInProgress, JobStatusPaused, JobStatusFailed:
		return true
	}
	return false
}

func (s *JobStatus) UnmarshalJSON(b []byte) error {
	var val string
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	tmp := JobStatus(val)
	if !tmp.IsValid() {
		return fmt.Errorf("invalid job status: %s", val)
	}
	*s = tmp
	return nil
}

// Job represent a crawl job are dispatched regularly.
type Job struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	URL            string    `grom:"not null;uniqueIndex"`
	RetryAttempts  int       `gorm:"default:0;check:retry_attempts >= 0"`
	Status         JobStatus `gorm:"type:varchar(20);not null"`
	Interval       time.Duration
	PauseRequested bool
	LastCrawledAt  *time.Time
	DispatchedAt   *time.Time
	NextRunAt      time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (j *Job) IsDue() bool {
	return j.NextRunAt.After(time.Now())
}

// ScheduleNextRun sets the NextRunAt and Status for the job.
// If the job has never been crawled (LastCrawledAt is nil),
// it schedules the job to run immediately (now).
// Otherwise, it schedules the next run as the later of:
//   - LastCrawledAt + Interval
//   - the current time (to avoid scheduling in the past).
func (j *Job) ScheduleNextRun() {
	var nextRun time.Time

	if j.LastCrawledAt != nil {
		nextRun = j.LastCrawledAt.Add(j.Interval)
		// ensure not scheduling in the past
		if nextRun.Before(time.Now()) {
			nextRun = time.Now()
		}
	} else {
		nextRun = time.Now()
	}

	j.NextRunAt = nextRun
	j.Status = JobStatusScheduled
}

// Pause pauses the job if the current state allows it.
//
// Behavior by job status:
//   - JobStatusInProgress: marks the job as pause-requested.
//   - JobStatusFailed: returns ErrCannotPause.
//   - Any other state: immediately sets the job status to JobStatusPaused.
func (j *Job) Pause() error {
	switch j.Status {
	case JobStatusFailed:
		return ErrCannotPause
	case JobStatusInProgress:
		j.PauseRequested = true
		return nil
	default:
		j.Status = JobStatusPaused
		return nil
	}
}

// Resume restarts a job if it was paused or failed.
//
//   - If the job is paused or failed, RetryAttempts is reset to 0,
//     the next run is scheduled, and the status becomes Scheduled.
//   - Jobs that are already scheduled or in progress are not changed.
func (j *Job) Resume() error {
	if j.Status == JobStatusScheduled || j.Status == JobStatusInProgress {
		return nil
	}
	j.RetryAttempts = 0
	j.ScheduleNextRun()
	return nil
}

// UpdateInterval changes the interval and schedules next run
func (j *Job) UpdateInterval(interval time.Duration) {
	j.Interval = interval
	j.ScheduleNextRun()
}

func (j *Job) ShouldRetry(maxAttempts int) bool {
	return j.RetryAttempts < maxAttempts
}

type JobDispatched struct {
	ID           uint
	URL          string
	DispatchedAt time.Time
}
