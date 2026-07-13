package responses

type PaginationResponse[T any] struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	Items      []T   `json:"items"`
}
