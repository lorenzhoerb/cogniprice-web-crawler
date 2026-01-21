package dispatcher

import (
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/shared/logger"
	"go.uber.org/zap"
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
		logger.Log.Info("dispatching job",
			// structured fields
			zap.Uint64("jobID", uint64(job.ID)),
			zap.String("url", job.URL),
		)
	}
	return nil
}
