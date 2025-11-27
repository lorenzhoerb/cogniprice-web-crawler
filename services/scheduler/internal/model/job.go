package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// JobStatus represents the current status of a job.
//
// Status values:
//   - Scheduled: The job is ready to be dispatched if NextRunAt <= now.
//   - Running: The job was dispatched successfully and is currently executing.
//   - Paused: The job will not be scheduled until resumed.
//   - Failed: The job exceeded its failure/retry limit and is no longer eligible for dispatch.
type JobStatus int

const (
	Scheduled = iota
	InProgress
	Paused
	Failed
)

func (j JobStatus) String() string {
	switch j {
	case Scheduled:
		return "SCHEDULED"
	case InProgress:
		return "IN_PROGRESS"
	case Paused:
		return "PAUSED"
	case Failed:
		return "FAILED"
	default:
		return "UNKNOWN"
	}
}

func (j JobStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(j.String())
}

func (j *JobStatus) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch s {
	case "scheduled":
		*j = Scheduled
	case "inProgress":
		*j = InProgress
	case "paused":
		*j = Paused
	case "failed":
		*j = Failed
	default:
		return fmt.Errorf("unknown JobStatus: %s", s)
	}
	return nil
}

var (
	ErrCannotPause = errors.New("cannot pause job in current state")
)

// Job represent a crawl job are dispatched regularly.
type Job struct {
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	URL            string    `grom:"not null;uniqueIndex"`
	RetryAttempts  int       `gorm:"default:0;check:retry_attempts >= 0"`
	Status         JobStatus `gorm:"type:int;not null"`
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
	j.Status = Scheduled
}

// Pause sets the job status to Paused if allowed
func (j *Job) Pause() error {
	if j.Status == Failed {
		return ErrCannotPause
	}
	j.Status = Paused
	return nil
}

// Resume sets the job status to Scheduled if it was paused
func (j *Job) Resume() error {
	if j.Status == Scheduled || j.Status == InProgress {
		return nil
	}
	j.Status = Scheduled
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
