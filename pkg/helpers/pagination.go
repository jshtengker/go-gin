package helpers

type PaginationParams struct {
	Page     int
	PageSize int
	Offset   int
}

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

func NewPaginationParams(page, pageSize int) PaginationParams {
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}
