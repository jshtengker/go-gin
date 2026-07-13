package helpers

import (
	"backend/internal/api/responses"
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

func ParsePagination(c *gin.Context) PaginationParams {
	page, err := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	if err != nil || page < 1 {
		page = DefaultPage
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)))
	if err != nil || pageSize < 1 {
		pageSize = MaxPageSize
	}

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
	}
}

func BuildPagination[T any](params PaginationParams, total int64, items []T) responses.PaginationResponse[T] {
	totalPages := int(math.Ceil(float64(total) / float64(params.PageSize)))

	if totalPages == 0 {
		totalPages = 1
	}

	return responses.PaginationResponse[T]{
		Total:      total,
		Page:       params.Page,
		PageSize:   params.PageSize,
		TotalPages: totalPages,
		Items:      items,
	}
}
