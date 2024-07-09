package utils

import (
	"math"
	"strconv"

	"net/http"
)

type PaginationParams struct {
	Page       int
	Limit      int
	Search     string
	SortBy     string
	SortOrder  string
	Category   string
	BorrowerId int
	BorrowType string
	BorrowDate string
}

const (
	maxLimit      = 100.0
	pageKey       = "pageNumber"
	limitKey      = "itemsPerPage"
	searchKey     = "search"
	sortByKey     = "sortBy"
	sortOrderKey  = "sortOrder"
	categoryKey   = "category"
	borrowerId    = "readerId"
	borrowTypeKey = "borrowStatus"
	borrowDateKey = "borrowDate"
)

func parsePage(r *http.Request) int {
	pageStr := r.URL.Query().Get(pageKey)
	page, _ := strconv.ParseInt(pageStr, 10, 32)
	page = int64(math.Max(1.0, float64(page)))
	return int(page)
}

func parseLimit(r *http.Request) int {
	limitStr := r.URL.Query().Get(limitKey)
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	limit = int64(math.Max(1.0, math.Min(maxLimit, float64(limit))))
	return int(limit)
}

func parseBorrowerId(r *http.Request) int {
	limitStr := r.URL.Query().Get(borrowerId)
	limit, _ := strconv.ParseInt(limitStr, 10, 32)
	limit = int64(math.Max(0.0, math.Min(maxLimit, float64(limit))))
	return int(limit)
}

func CountTotalPages(limit, totalItems int) int {
	return int(math.Ceil(float64(totalItems) / math.Max(1.0, float64(limit))))
}

func GetPaginationParams(r *http.Request, defaultSortBy, defaultSortOrder string) PaginationParams {
	params := PaginationParams{
		Page:       1,
		Limit:      6,
		Search:     "",
		SortBy:     defaultSortBy,
		SortOrder:  defaultSortOrder,
		Category:   "",
		BorrowerId: -1,
		BorrowType: "",
		BorrowDate: "",
	}

	for k := range r.URL.Query() {
		switch k {
		case pageKey:
			// parse page number
			params.Page = parsePage(r)

		case limitKey:
			// parse limit
			params.Limit = parseLimit(r)

		case searchKey:
			// parse search term
			params.Search = r.URL.Query().Get(searchKey)

		case sortByKey:
			// parse sort by
			params.SortBy = r.URL.Query().Get(sortByKey)

		case sortOrderKey:
			// parse sort order
			params.SortOrder = r.URL.Query().Get(sortOrderKey)

		case categoryKey:
			// category type
			params.Category = r.URL.Query().Get(categoryKey)

		case borrowerId:
			// category type
			params.BorrowerId = parseBorrowerId(r)

		case borrowTypeKey:
			//borrow status type
			params.BorrowType = r.URL.Query().Get(borrowTypeKey)

		case borrowDateKey:
			//borrow date
			params.BorrowDate = r.URL.Query().Get(borrowDateKey)
		}

	}

	return params
}

func GetSortingData(r *http.Request, defaultSortBy, defaultSortOrder string) (sortBy, sortOrder string) {
	params := GetPaginationParams(r, defaultSortBy, defaultSortOrder)
	return params.SortBy, params.SortOrder
}
