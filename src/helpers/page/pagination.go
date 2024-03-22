package page

import (
	"net/http"
	"strconv"
)

const (
	MaxPageSize     = 100
	DefaultPageSize = 10
)

func GetPage(r *http.Request) int {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		return 1
	}
	return page
}

func GetPageSize(r *http.Request) int {
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	if pageSize < 0 {
		pageSize = -1
	}
	if pageSize == 0 {
		return DefaultPageSize
	}
	if pageSize > MaxPageSize {
		return MaxPageSize
	}
	return pageSize
}

func GetPageOffset(page, pageSize int) (offset int) {
	offset = 0
	if page > 0 {
		offset = (page - 1) * pageSize
	}
	return
}

func Limit(page, pageSize int) (limit, offset int) {
	// 不进行分页
	if pageSize == -1 {
		return -1, -1
	}
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}
	limit = pageSize
	offset = 0
	if page > 0 {
		offset = (page - 1) * pageSize
	}
	return
}
