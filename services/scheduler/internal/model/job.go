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
	DispatchedAt   *time.Time
	NextRunAt      time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}

func (j *Job) IsDue() bool {
	return j.NextRunAt.After(time.Now())
}

// ScheduleNextRun updates NextRunAt based on Interval
func (j *Job) ScheduleNextRun() {
	j.NextRunAt = time.Now().Add(j.Interval)
	j.Status = JobStatusScheduled
}

// Pause sets the job status to Paused if allowed
func (j *Job) Pause() error {
	if j.Status == JobStatusFailed {
		return ErrCannotPause
	}
	j.Status = JobStatusPaused
	return nil
}

// Resume sets the job status to Scheduled if it was paused
func (j *Job) Resume() error {
	if j.Status == JobStatusScheduled || j.Status == JobStatusInProgress {
		return nil
	}
	j.Status = JobStatusScheduled
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
