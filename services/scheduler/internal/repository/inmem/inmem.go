package inmem

import (
	"sync"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
)

type inmemJobRepository struct {
	data map[uint]*model.Job
	mu   sync.RWMutex
}

func New() *inmemJobRepository {
	return &inmemJobRepository{
		data: map[uint]*model.Job{},
	}
}

func (r *inmemJobRepository) GetJob(id uint) (*model.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.data[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return v, nil
}

func (r *inmemJobRepository) PutJob(job *model.Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[job.ID] = job
	return nil
}

func (r *inmemJobRepository) UpdateJobs(jobs []*model.Job) error {
	return nil
}

func (r *inmemJobRepository) GetDueJobs(limit int) ([]*model.Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := []*model.Job{}
	for _, job := range r.data {
		if job.IsDue() {
			jobs = append(jobs, job)
			if len(jobs) >= limit {
				break
			}
		}
	}
	return jobs, nil
}
