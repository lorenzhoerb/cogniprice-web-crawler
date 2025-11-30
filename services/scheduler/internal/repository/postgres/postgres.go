package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/model"
	"github.com/lorenzhoerb/cogniprice/services/scheduler/internal/repository"
	"github.com/lorenzhoerb/cogniprice/shared/pagination"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *jobRepository) SaveAll(jobs []*model.Job) error {
	return r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&jobs).Error
}

// Get due Jobs ,
func (r *jobRepository) GetDue(limit int) ([]*model.Job, error) {
	var jobs []*model.Job
	db := r.db.
		Where("next_run_at <= ?", time.Now()).
		Where("status = ? ", model.JobStatusScheduled).
		Order("next_run_at ASC")

	if limit > 0 {
		db = db.Limit(limit)
	}

	result := db.Find(&jobs)
	if result.Error != nil {
		return nil, result.Error
	}

	return jobs, nil
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

func (r *jobRepository) List(filter *model.ListJobsFilter) ([]*model.Job, *pagination.Pagination, error) {
	var jobs []*model.Job

	sortBy, sortOrder := sanitizeSort(filter.SortBy, filter.SortOrder)
	pagination := pagination.NewPagination(filter.Page, filter.PageSize)

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
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
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

// sanitizeSort safely returns the column and direction for sorting.
// It uses defaults if nil or invalid values are provided and prevents SQL injection.
func sanitizeSort(sortBy, orderBy *string) (column, direction string) {
	// Define allowed values
	allowedColumns := map[string]bool{
		"status":      true,
		"url":         true,
		"created_at":  true,
		"next_run_at": true,
	}
	allowedDirections := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	// Default values
	column = "created_at"
	direction = "asc"

	if sortBy != nil && allowedColumns[*sortBy] {
		column = *sortBy
	}

	if orderBy != nil && allowedDirections[*orderBy] {
		direction = *orderBy
	}

	return
}
