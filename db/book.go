package db

import (
	"fmt"
	"librarymanagement/logger"
	"librarymanagement/web/utils"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Book struct {
	Id         int       `json:"id"`
	Title      string    `db:"title"       validate:"required"   json:"title"`
	Category   string    `db:"category"    validate:"required"   json:"category"`
	Author     string    `db:"author"      validate:"required"   json:"author"`
	Quantity   int       `db:"quantity"    validate:"required"   json:"quantity"`
	Available  int       `db:"available"   validate:"required"   json:"available"`
	Summary    string    `db:"summary"     validate:"required"   json:"summary"`
	TotalPage  int       `db:"total_page"  validate:"required"   json:"total_page"`
	ImageLink  string    `db:"image_link"  validate:"required"   json:"image_link"`
	Created_at time.Time `db:"created_at"`
}

type BookRepo struct {
	Table string
}

var bookRepo *BookRepo

func InitBookRepo() {
	bookRepo = &BookRepo{Table: "book"}
}

func GetBookRepo() *BookRepo {
	return bookRepo
}

func (r *BookRepo) InsertBook(book *Book) (*Book, error) {

	column := map[string]interface{}{
		"title":      book.Title,
		"category":   book.Category,
		"author":     book.Author,
		"quantity":   book.Quantity,
		"available":  book.Available,
		"summary":    book.Summary,
		"total_page": book.TotalPage,
		"image_link": book.ImageLink,
	}

	var columns []string
	var values []any
	for columnName, columnValue := range column {
		columns = append(columns, columnName)
		values = append(values, columnValue)
	}
	qry, args, err := GetQueryBuilder().
		Insert(r.Table).
		Columns(columns...).
		Suffix(`
			RETURNING 		
			title,
			category,
			author,
			quantity,
			available,
			summary,
			total_page,
			image_link
		`).
		Values(values...).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new book insert query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	var insBook Book
	err = GetReadDB().QueryRow(qry, args...).Scan(&insBook.Title, &insBook.Category, &insBook.Author, &insBook.Quantity, &insBook.Available, &insBook.Summary, &insBook.TotalPage, &insBook.ImageLink)
	if err != nil {
		slog.Error(
			"Failed to execute insert book query",
			logger.Extra(map[string]interface{}{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	return &insBook, nil

}

func (r *BookRepo) UpdateBook(uBook *Book) error {

	updateQry, args, err := GetQueryBuilder().Update(r.Table).
		Set("title", uBook.Title).
		Set("author", uBook.Author).
		Set("category", uBook.Category).
		Set("quantity", uBook.Quantity).
		Set("available", uBook.Available).
		Set("summary", uBook.Summary).
		Set("total_page", uBook.TotalPage).
		Set("image_link", uBook.ImageLink).
		Where(sq.Eq{"Id": uBook.Id}).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create update query of book",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": updateQry,
				"args":  args,
			}),
		)
		return err
	}

	_, err = GetReadDB().Exec(updateQry, args...)
	if err != nil {
		slog.Error(
			"Failed to update book info",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": updateQry,
				"args":  args,
			}),
		)
		return err
	}

	return nil
}

func (r *BookRepo) DeleteBook(id int) error {

	delQry, args, err := GetQueryBuilder().Delete(r.Table).Where(sq.Eq{"id": id}).ToSql()
	if err != nil {
		slog.Error(
			"Failed to create delete query of book",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": delQry,
				"args":  args,
			}),
		)
		return err
	}

	_, err = GetReadDB().Exec(delQry, args...)
	if err != nil {
		slog.Error(
			"Failed to delete book",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": delQry,
				"args":  args,
			}),
		)
		return err
	}

	return nil

}

func (r *BookRepo) GetBookList(params utils.PaginationParams) (utils.Page, error) {

	//Start quering book details
	type BookResult struct {
		items []*Book
		err   error
	}

	itemsChan := make(chan BookResult)

	go func() {
		books, err := r.GetBooks(params)
		if err != nil {
			itemsChan <- BookResult{
				err: err,
			}
			return
		}

		itemsChan <- BookResult{
			items: books,
		}
	}()

	// Start quering book count
	type CntResult struct {
		totalCnts int
		err       error
	}

	totalCntChan := make(chan CntResult)

	go func() {
		totalCount, err := r.GetTotal(params)
		if err != nil {
			totalCntChan <- CntResult{
				err: err,
			}
			return
		}

		totalCntChan <- CntResult{
			totalCnts: totalCount,
		}
	}()

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
		case result := <-itemsChan:
			if result.err != nil {
				slog.Error(
					"Failed to get Books",
					logger.Extra(map[string]any{
						"error": result.err.Error(),
					}),
				)
				return page, result.err
			}
			page.Items = result.items

		case result := <-totalCntChan:
			if result.err != nil {
				slog.Error(
					"Failed to get audit logs count",
					logger.Extra(map[string]any{
						"error": result.err.Error(),
					}),
				)
				return page, result.err
			}
			page.TotalItems = result.totalCnts
		}
	}

	page.TotalPages = utils.CountTotalPages(params.Limit, page.TotalItems)

	return page, nil

}

func (r *BookRepo) GetBooks(params utils.PaginationParams) ([]*Book, error) {

	query := r.BuildFilterQuery(params)
	queryString, args, err := query.ToSql()
	if err != nil {
		slog.Error(
			"Failed to create books select query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": queryString,
				"args":  args,
			}),
		)
		return nil, err
	}

	books := []*Book{}
	err = GetReadDB().Select(&books, queryString, args...)
	if err != nil {
		slog.Error(
			"Failed to get books",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return nil, err
	}
	return books, nil
}

func (r *BookRepo) BuildFilterQuery(params utils.PaginationParams) sq.SelectBuilder {
	query := GetQueryBuilder().
		Select("*").
		From(r.Table)

	if params.Category != "" {
		query = query.Where(sq.Eq{"category": params.Category})
	}

	if params.Search != "" {
		likePattern := fmt.Sprintf("%%%s%%", params.Search)
		query = query.Where(
			sq.Or{
				sq.Like{"title": likePattern},
				sq.Like{"author": likePattern},
			},
		)
	}

	query = query.Limit(uint64(params.Limit)).
		Offset(uint64((params.Page - 1) * params.Limit)).
		OrderBy(params.SortBy + " " + params.SortOrder)

	return query
}

func (r *BookRepo) GetTotal(params utils.PaginationParams) (int, error) {
	query := r.BuildCntQuery(params)
	queryString, args, err := query.ToSql()
	if err != nil {
		slog.Error(
			"Failed to create books count select query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": queryString,
				"args":  args,
			}),
		)
		return 0, err
	}

	var totalCount int
	err = GetWriteDB().Get(&totalCount, queryString, args...)
	if err != nil {
		slog.Error(
			"Failed to get book count",
			logger.Extra(map[string]any{
				"error": err.Error(),
			}),
		)
		return 0, err
	}

	return totalCount, nil
}

func (r *BookRepo) BuildCntQuery(params utils.PaginationParams) sq.SelectBuilder {
	query := GetQueryBuilder().
		Select("COUNT(created_at)").
		From(r.Table)

	if params.Category != "" {
		query = query.Where(sq.Eq{"category": params.Category})
	}

	if params.Search != "" {
		likePattern := fmt.Sprintf("%%%s%%", params.Search)
		query = query.Where(
			sq.Or{
				sq.Like{"title": likePattern},
				sq.Like{"author": likePattern},
			},
		)
	}

	return query
}

func (r *BookRepo) UpdateBookAvailability(bookId int) error {

	updateQry, args, err := GetQueryBuilder().Update(r.Table).
		Set("available", sq.Expr("available - 1")).
		Where(sq.Eq{"id": bookId}).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create update for book table decrease available field by one",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": updateQry,
				"args":  args,
			}),
		)
		return err
	}

	// Execute the update query
	_, err = GetReadDB().Exec(updateQry, args...)
	if err != nil {
		slog.Error(
			"Failed to update of book count decrease by one",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": updateQry,
				"args":  args,
			}),
		)
		return err
	}

	return nil
}

func (r *BookRepo) UpdateBookCount(bookId int) error {

	updateQry, args, err := GetQueryBuilder().Update(r.Table).
		Set("available", sq.Expr("available + 1")).
		Where(sq.Eq{"id": bookId}).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create update for book table increase available field by one",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": updateQry,
				"args":  args,
			}),
		)
		return err
	}

	// Execute the update query
	_, err = GetReadDB().Exec(updateQry, args...)
	if err != nil {
		slog.Error(
			"Failed to update of book count decrease by one",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": updateQry,
				"args":  args,
			}),
		)
		return err
	}

	return nil
}
