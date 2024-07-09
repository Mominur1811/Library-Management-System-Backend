package db

import (
	"fmt"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type BorrowRequest struct {
	RequestId    *int       `db:"request_id"`
	BookId       int        `db:"book_id"          json:"book_id"`
	BorrowerId   int        `db:"borrower_id"        json:"borrower_id"`
	IssuedAt     *time.Time `db:"issued_at"        json:"issued_at"`
	ReturnedAt   *time.Time `db:"returned_at"      json:"returned_at"`
	ReadPage     int        `db:"read_page"        json:"read_page"`
	BorrowStatus string     `db:"borrow_status"   json:"borrow_status"`
}

type BorrowHistory struct {
	RequestId     *int       `db:"request_id"         json:"request_id"`
	BookId        int        `db:"book_id"            json:"book_id"`
	BorrowerId    int        `db:"borrower_id"        json:"borrower_id"`
	BorrowerName  string     `db:"borrower_name"      json:"borrower_name"`
	BookTitle     string     `db:"book_title"         json:"book_title"`
	BookAvailable int        `db:"book_available"     json:"book_available"`
	ReadPage      int        `db:"read_page"          json:"read_page"`
	Totalpage     int        `db:"total_page"         json:"total_page"`
	BorrowStatus  string     `db:"borrow_status"      json:"borrow_status"`
	IssuedAt      *time.Time `db:"issued_at"          json:"issued_at"`
	ReturnedAt    *time.Time `db:"returned_at"        json:"returned_at"`
	Requestdate   *time.Time `db:"created_at"         json:"created_at"`
}

// Push borrow request
func (r *BookRequestRepo) PushBorrowRequest(requestBook *BorrowRequest) (*BorrowRequest, error) {

	column := map[string]interface{}{
		"book_id":       requestBook.BookId,
		"borrower_id":   requestBook.BorrowerId,
		"borrow_status": requestBook.BorrowStatus,
	}

	var columns []string
	var values []any
	for columnName, columnValue := range column {
		columns = append(columns, columnName)
		values = append(values, columnValue)
	}
	qry, args, err := GetQueryBuilder().
		Insert("borrow_history").
		Columns(columns...).
		Suffix(`
			RETURNING
			request_id, 		
			book_id,
			borrower_id,
			issued_at,
			returned_at,
			read_page,
			borrow_status
		`).
		Values(values...).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new borrow request query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	fmt.Println(qry)
	var insRequest BorrowRequest
	err = GetReadDB().QueryRow(qry, args...).Scan(&insRequest.RequestId, &insRequest.BookId, &insRequest.BorrowerId, &insRequest.IssuedAt, &insRequest.ReturnedAt, &insRequest.ReadPage, &insRequest.BorrowStatus)
	if err != nil {
		slog.Error(
			"Failed to execute book request query",
			logger.Extra(map[string]interface{}{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}
	return &insRequest, nil
}

// Borrow History (Approved, Pending, Request, Returned)
func (r *BookRequestRepo) GetBorrowStatus(params utils.PaginationParams) (utils.Page, error) {

	// Fetch Borrow History
	type BorrowResult struct {
		BorrowHistoy []*BorrowHistory
		err          error
	}
	borrowResultChan := make(chan BorrowResult)
	go func() {

		borrowHistory, err := r.GetBorrowHistory(params)
		if err != nil {
			borrowResultChan <- BorrowResult{
				err: err,
			}
			return
		}

		borrowResultChan <- BorrowResult{
			BorrowHistoy: borrowHistory,
		}
	}()

	//Fetch Total Number of Row Selected
	type BorrowResultCount struct {
		count int
		err   error
	}
	borrowResultCountChan := make(chan BorrowResultCount)

	go func() {

		count, err := r.GetBorrowResultCount(params)
		if err != nil {
			borrowResultCountChan <- BorrowResultCount{
				err: err,
			}
			return
		}
		borrowResultCountChan <- BorrowResultCount{
			count: *count,
		}

	}()

	//Create Page
	page := utils.Page{
		Items:        []any{},
		ItemsPerPage: params.Limit,
		PageNumber:   params.Page,
		TotalItems:   0,
		TotalPages:   0,
	}

	// wait for page data and total items
	for i := 1; i <= 2; i++ {
		select {
		case result := <-borrowResultChan:
			if result.err != nil {
				slog.Error(
					"Failed to get borrow history",
					logger.Extra(map[string]any{
						"error": result.err.Error(),
					}),
				)
				return page, result.err
			}
			page.Items = result.BorrowHistoy

		case result := <-borrowResultCountChan:
			if result.err != nil {
				slog.Error(
					"Failed to get total history count",
					logger.Extra(map[string]any{
						"error": result.err.Error(),
					}),
				)
				return page, result.err
			}
			page.TotalItems = result.count
		}
	}

	page.TotalPages = utils.CountTotalPages(params.Limit, page.TotalItems)

	return page, nil

}

func (r *BookRequestRepo) GetBorrowHistory(params utils.PaginationParams) ([]*BorrowHistory, error) {

	query := GetQueryBuilder().Select("rq.request_id", "rq.book_id", "rq.borrower_id", "r.name AS borrower_name", "b.title AS book_title",
		"b.available AS book_available", "rq.read_page", "b.total_page AS total_page", "rq.borrow_status", "rq.issued_at",
		"rq.returned_at", "rq.created_at").
		From("borrow_history rq").
		Join("book b ON rq.book_id = b.id").
		Join("reader r ON rq.borrower_id = r.id")

	if params.Search != "" {
		query = query.Where(
			sq.Or{
				sq.Expr("LOWER(b.title) LIKE ?", "%"+strings.ToLower(params.Search)+"%"),
				sq.Expr("LOWER(b.author) LIKE ?", "%"+strings.ToLower(params.Search)+"%"),
			},
		)
	}

	if params.Category != "" {
		query = query.Where(sq.Eq{"b.category": params.Category})
	}

	if params.BorrowerId != -1 {
		query = query.Where(sq.Eq{"r.id": params.BorrowerId})
	}

	if params.BorrowType != "" {
		query = query.Where(sq.Eq{"rq.borrow_status": params.BorrowType})
	}

	query = query.Limit(uint64(params.Limit)).
		Offset(uint64((params.Page - 1) * params.Limit)).
		OrderBy(params.SortBy + " " + params.SortOrder)

	if params.BorrowDate != "" {

		borrowDate, err := time.Parse("2006-01-02", params.BorrowDate)
		if err != nil {
			return nil, err
		}

		if params.BorrowType == "Approved" {
			query = query.Where(sq.Gt{"rq.issued_at": borrowDate})
		} else {
			query = query.Where(sq.Gt{"rq.created_at": params.BorrowDate})
		}

	}

	qry, args, err := query.ToSql()
	fmt.Println(qry, args)

	if err != nil {
		slog.Error(
			"Failed to borrow history query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}
	borrowHistory := []*BorrowHistory{}
	err = GetReadDB().Select(&borrowHistory, qry, args...)
	if err != nil {
		slog.Error(
			"Failed to Fetch borrow history",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	return borrowHistory, nil
}

func (r *BookRequestRepo) GetBorrowResultCount(params utils.PaginationParams) (*int, error) {

	query := GetQueryBuilder().Select("Count(*)").
		From("borrow_history rq").
		Join("book b ON rq.book_id = b.id").
		Join("reader r ON rq.borrower_id = r.id")

	if params.Search != "" {
		likePattern := fmt.Sprintf("%%%s%%", params.Search)
		query = query.Where(
			sq.Or{
				sq.Like{"b.title": likePattern},
				sq.Like{"b.author": likePattern},
			},
		)
	}

	if params.Category != "" {
		query = query.Where(sq.Eq{"b.category": params.Category})
	}

	if params.BorrowerId != -1 {
		query = query.Where(sq.Eq{"r.id": params.BorrowerId})
	}

	if params.BorrowType != "" {
		query = query.Where(sq.Eq{"rq.borrow_status": params.BorrowType})
	}

	if params.BorrowDate != "" {

		borrowDate, err := time.Parse("2006-01-02", params.BorrowDate)
		if err != nil {
			return nil, err
		}

		if params.BorrowType == "Approved" {
			query = query.Where(sq.Gt{"rq.issued_at": borrowDate})
		} else {
			query = query.Where(sq.Gt{"rq.created_at": params.BorrowDate})
		}

	}

	qry, args, err := query.ToSql()
	if err != nil {
		slog.Error(
			"Failed to create Get Borrow Table History Count Query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	var totalCount int
	err = GetWriteDB().Get(&totalCount, qry, args...)
	if err != nil {
		slog.Error(
			"Failed to Fetch upapproved user",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	return &totalCount, nil
}
