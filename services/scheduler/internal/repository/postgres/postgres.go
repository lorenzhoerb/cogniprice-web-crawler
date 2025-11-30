package postgres

import (
	"errors"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
	"github.com/lorenzhoerb/cogniprice/shared"
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
	result := r.db.Where("url = ?", url).Take(&job)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, repository.ErrNotFound // Not found, return nil without error
		}
		return nil, result.Error
	}
	return &job, nil
}

func (r *jobRepository) List(filter *model.ListJobsFilter) ([]*model.Job, *shared.Pagination, error) {
	var jobs []*model.Job

	pagination := shared.NewPagination(filter.Page, filter.PageSize)

	db := r.db.Model(&model.Job{})

	// Apply URL filter if provided
	if filter.URL != nil {
		db = db.Where("url ILIKE ?", "%"+*filter.URL+"%")
	}

	// Apply Status filter if provided
	if filter.Status != nil {
		db = db.Where("status = ?", filter.Status)
	}

	// Get total count (ignoring limit/offset)
	if err := db.Count(&pagination.Total).Error; err != nil {
		return nil, nil, err
	}

	// Apply pagination
	result := db.
		Limit(pagination.Limit()).
		Offset(pagination.Offset()).
		Find(&jobs)

	if result.Error != nil {
		return nil, nil, result.Error
	}

	return jobs, pagination, nil
}

func (r *jobRepository) Delete(id int) error {
	return r.db.Delete(&model.Job{}, id).Error
}
