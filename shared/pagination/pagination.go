package pagination

import "math"

const (
	DEFAULT_PAGE_SIZE int = 25
	MAX_PAGE_SIZE     int = 100
)

// Pagination represents a pagination data structure.
type Pagination struct {
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
	Total    int64 `json:"total"`
}

// Creates a new Pagination instance
func NewPagination(page, pageSize int) *Pagination {
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = DEFAULT_PAGE_SIZE
	} else if page >= MAX_PAGE_SIZE {
		pageSize = MAX_PAGE_SIZE
	}

	return &Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// Limit return the SQL Limit
func (p *Pagination) Limit() int {
	return p.PageSize
}

// Offset returns the SQL offset
func (p *Pagination) Offset() int {
	return (1 - p.Page) * p.PageSize
}

// TotalPages calculates the number of pages
func (p *Pagination) SetTotal(total int64) {
	p.Total = total
}

func (p *Pagination) TotalPages() int {
	pages := int(math.Ceil(float64(p.Total) / float64(p.PageSize)))
	if pages == 0 {
		pages = 1
	}
	return pages
}

func (p *Pagination) CurrentPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}
