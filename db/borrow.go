package db

import (
	"fmt"
	"librarymanagement/logger"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Request struct {
	BookTitle      string     `db:"book_title"       json:"book_title"`
	BookAvailable  int        `db:"book_available"   json:"book_available"`
	ReaderUsername string     `db:"reader_username"  json:"reader_name"`
	RequestId      *int       `db:"request_id"`
	BookID         int        `db:"bookid"           json:"book_id"`
	ReaderID       int        `db:"readerid"         json:"reader_id"`
	IssuedAt       *time.Time `db:"issued_at"        json:"issued_at"`
	RequestStatus  string     `db:"request_status"   json:"request_status"`
}


type BookRequestRepo struct {
	Table string
}

var bookRequestRepo *BookRequestRepo

func InitBookReqeustRepo() {
	bookRequestRepo = &BookRequestRepo{Table: "book_request"}
}

func GetBookRequestRepo() *BookRequestRepo {
	return bookRequestRepo
}


func (r *BookRequestRepo) ChangeBorrowStatus(reqId int, updatedStatus string) error {

	updateQry, args, err := GetQueryBuilder().Update("borrow_history").
		Set("borrow_status", updatedStatus).
		Set("issued_at", time.Now()).
		Where(sq.Eq{"request_id": reqId}).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create update query for request status",
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
			"Failed to update of readers borrow status",
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


func (r *BookRequestRepo) GetBorrowedBooks() ([]*Request, error) {

	qry, args, err := GetQueryBuilder().Select("b.title AS book_title", "b.available AS book_available", "r.name AS reader_username", "rq.request_id", "rq.bookid", "rq.readerid", "rq.issued_at", "rq.request_status").
		From("book_request rq").
		Join("book b ON rq.bookid = b.id").
		Join("reader r ON rq.readerid = r.id").
		Where(sq.Eq{"rq.request_status": "Approved"}).ToSql()

	if err != nil {
		slog.Error(
			"Failed to create Get Unapproved request query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	borrowedBooks := []*Request{}
	err = GetReadDB().Select(&borrowedBooks, qry, args...)
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

	return borrowedBooks, nil
}

func (r *BookRequestRepo) UpdateUserReadProgress(reqId int, read_today int) error {

	updateQry, args, err := GetQueryBuilder().Update("borrow_history").
		Set("read_page", sq.Expr(fmt.Sprintf("read_page + %d", read_today))).
		Where(sq.Eq{"request_id": reqId}).
		ToSql()
	fmt.Println(updateQry)
	if err != nil {
		slog.Error(
			"Failed to create update query for request status",
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
			"Failed to update of readers active status",
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

func (r *BookRequestRepo) ValidateUserBorrowRequest(req *BorrowRequest) error {

	fmt.Println(req)
	qry, args, err := GetQueryBuilder().Select("COUNT(*)").
		From("borrow_history").
		Where(sq.Eq{"borrower_id": req.BorrowerId}).Where(sq.Eq{"book_id": req.BookId}).
		Where(sq.Or{
			sq.Eq{"borrow_status": "Approved"},
			sq.Eq{"borrow_status": "Pending"}}).ToSql()

	if err != nil {
		slog.Error(
			"Failed to create Get Unapproved request query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil
	}

	var count int
	err = GetReadDB().Get(&count, qry, args...)
	if err != nil {
		slog.Error(
			"Failed to Fetch upapproved user",
			logger.Extra(map[string]any{
				"error":   err.Error(),
				"payload": count,
			}),
		)
		return err
	}

	if count != 0 {
		err = fmt.Errorf("already hold books and pending book request")
		slog.Error(
			"User has already Hold or pending book request",
			logger.Extra(map[string]any{
				"error":   err.Error(),
				"payload": count,
			}),
		)
		return err
	}

	return nil
}
