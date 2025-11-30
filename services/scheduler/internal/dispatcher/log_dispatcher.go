package dispatcher

import (
	"fmt"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
)

// logDispatcher is a lightweight mock dispatcher.
// It simulates dispatching jobs by printing them to stdout instead of sending them
// to an actual queue or downstream service.
type logDispatcher struct{}

// NewLogDispatcher returns a new instance of LogDispatcher.
// This implementation is primarily intended for development, testing, and debugging,
// where observing dispatched jobs via logs is sufficient.
func NewLogDispatcher() *logDispatcher {
	return &logDispatcher{}
}

func (d *logDispatcher) DispatchJobs(jobs []model.JobDispatched) error {
	for _, job := range jobs {
		fmt.Printf("dispatching job: id=%d, url=%s\n", uint64(job.ID), job.URL)
	}
	return nil
}
