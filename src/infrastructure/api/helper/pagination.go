package helper

import (
	"github.com/hzmat24/api/infrastructure/api/dto"
	"strconv"
)

const DefaultPageLimit = 10

type IPaginationContext interface {
	QueryParam(key string) string
}

func NewPagination(c IPaginationContext, count int64) *dto.Pagination {
	limit := PaginationLimit(c)
	offset := PaginationOffset(c)

	return &dto.Pagination{
		Limit:      limit,
		Offset:     offset,
		HasPrev:    offset > 0,
		HasNext:    limit+offset < count,
		TotalCount: count,
	}
}

func PaginationLimit(c IPaginationContext) int64 {
	limit, err := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err == nil {
		return limit
	}

	return DefaultPageLimit
}

func PaginationOffset(c IPaginationContext) int64 {
	limit, err := strconv.ParseInt(c.QueryParam("offset"), 10, 64)
	if err == nil {
		return limit
	}

	return 0
}
