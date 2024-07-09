package db

import (
	"librarymanagement/logger"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type LoginReaderCredintials struct {
	Email    string `db:"email"    validate:"required,email"          json:"email"`
	Password string `db:"password" validate:"required"                json:"password"`
}

type Reader struct {
	Id        *int   `db:"id"       json:"id"`
	Name      string `db:"name"     validate:"required,alpha"          json:"name" `
	Email     string `db:"email"    validate:"required,email"          json:"email"`
	Password  string `db:"password" validate:"required"                json:"password"`
	Is_Active *bool  `db:"is_active"   `
}

type ReaderRepo struct {
	Table string
}

var readerRepo *ReaderRepo

func InitReaderRepo() {
	readerRepo = &ReaderRepo{Table: "reader"}
}

func GetReaderRepo() *ReaderRepo {
	return readerRepo
}

func (r *ReaderRepo) RegisterUser(newReader *Reader) (*Reader, error) {

	column := map[string]interface{}{
		"name":     newReader.Name,
		"password": newReader.Password,
		"email":    newReader.Email,
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
			name,
			password,
			email
		`).
		Values(values...).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new register query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	// Execute the SQL query and get the result
	var insertedReader Reader
	err = GetReadDB().QueryRow(qry, args...).Scan(&insertedReader.Name, &insertedReader.Password, &insertedReader.Email)
	if err != nil {
		slog.Error(
			"Failed to execute insert query",
			logger.Extra(map[string]interface{}{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	return &insertedReader, nil

}

func (r *ReaderRepo) GetUserInfo(userId int) (*Reader, error) {

	qry, args, err := GetQueryBuilder().Select("*").
		From("reader").
		Where(sq.Eq{"id": userId}).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create userInfo query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	var user Reader
	err = GetReadDB().Get(&user, qry, args...)
	if err != nil {
		slog.Error(
			"Failed to user info",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	return &user, nil

}

func (r *ReaderRepo) GetUnapprovedUser() ([]*Reader, error) {

	qry, args, err := GetQueryBuilder().Select("Id, Name", "Email", "Password").
		From("reader").
		Where(sq.Eq{"Is_Active": false}).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create Get Unapproved user query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	unapprovedReader := []*Reader{}
	err = GetReadDB().Select(&unapprovedReader, qry, args...)
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

	return unapprovedReader, nil

}

func (r *ReaderRepo) GetFetchUser() ([]*Reader, error) {

	qry, args, err := GetQueryBuilder().Select("Id, Name", "Email", "Password").
		From("reader").
		Where(sq.Eq{"Is_Active": true}).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create Get Unapproved user query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	unapprovedReader := []*Reader{}
	err = GetReadDB().Select(&unapprovedReader, qry, args...)
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

	return unapprovedReader, nil

}

func (r *ReaderRepo) ApprovedReader(email string) error {

	updateQry, args, err := GetQueryBuilder().Update(r.Table).
		Set("Is_Active", true).
		Where(sq.Eq{"Email": email}).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create update of readers active status query",
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

func (r *ReaderRepo) ValidateUserLogin(loginData LoginReaderCredintials) (*Reader, error) {

	qry, args, err := GetQueryBuilder().Select("*").From(`reader`).
		Where(sq.Eq{"email": loginData.Email, "password": loginData.Password, "is_active": true}).ToSql()
	if err != nil {
		slog.Error(
			"Failed to create validate login query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	// Execute the query
	var reader Reader
	err = GetReadDB().Get(&reader, qry, args...)
	if err != nil {
		slog.Error(
			"Failed to validate login",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}
	return &reader, nil
}
