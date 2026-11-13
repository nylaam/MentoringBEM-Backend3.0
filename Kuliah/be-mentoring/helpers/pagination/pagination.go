package pagination

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	TotalRows  int64 `json:"total_rows"`
	TotalPages int   `json:"total_pages"`
}

func New(page, limit int, totalRows int64) *Pagination {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	var totalPages int
	if totalRows > 0 {
		totalPages = int((totalRows + int64(limit) - 1) / int64(limit))
	} else {
		totalPages = 1
	}

	return &Pagination{
		Page:       page,
		Limit:      limit,
		TotalRows:  totalRows,
		TotalPages: totalPages,
	}
}

func (p *Pagination) Offset() int {
	return (p.Page - 1) * p.Limit
}
