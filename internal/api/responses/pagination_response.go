package responses

import (
	"backend/pkg/helpers"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaginationResponse[T any] struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
	Items      []T   `json:"items"`
}

func ParsePagination(c *gin.Context) helpers.PaginationParams {
	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(helpers.DefaultPage)))
	if err != nil {
		page = helpers.DefaultPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(helpers.DefaultPageSize)))
	if err != nil {
		pageSize = helpers.DefaultPageSize
	}

	return helpers.NewPaginationParams(page, pageSize)
}

func BuildPagination[T any](params helpers.PaginationParams, total int64, items []T) PaginationResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationResponse[T]{
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		Items:      items,
	}
}
