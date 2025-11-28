package postgres

import (
	"errors"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
	"gorm.io/gorm"
)

type jobRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *jobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Save(job *model.Job) error {
	return r.db.Save(job).Error
}

func (r *jobRepository) GetByID(id int) (*model.Job, error) {
	var job model.Job
	result := r.db.First(&job, id) // "id = ?" by default
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound // or return custom ErrNotFound
	}
	return &job, result.Error
}

func (r *jobRepository) GetByURL(url string) (*model.Job, error) {
	var job model.Job
	result := r.db.Where("url = ?", url).First(&job)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound // Not found, return nil without error
	}
	return &job, result.Error
}

func (r *jobRepository) Delete(id int) error {
	return r.db.Delete(&model.Job{}, id).Error
}
