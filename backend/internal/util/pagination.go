package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PageParam 分页参数。
type PageParam struct {
	Page     int
	PageSize int
}

// ParsePage 解析分页参数，默认 page=1、page_size=10，page_size 上限 100。
func ParsePage(c *gin.Context) PageParam {
	page := parseIntDefault(c.Query("page"), 1)
	pageSize := parseIntDefault(c.Query("page_size"), 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return PageParam{Page: page, PageSize: pageSize}
}

// Offset 计算偏移量。
func (p PageParam) Offset() int {
	return (p.Page - 1) * p.PageSize
}

func parseIntDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
