package db

import (
	"fmt"
	"librarymanagement/logger"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type Admin struct {
	Email    string `json:"email"      validate:"required,email"  db:"email"`
	Password string `json:"password"   validate:"required"        db:"password"`
}

type SubAdmin struct {
	Email string `json:"email"  db:"email"`
}

type AdminRepo struct {
	Table string
}

var adminRepo *AdminRepo

func InitAdminRepo() {
	adminRepo = &AdminRepo{Table: "admin"}
}

func GetAdminRepo() *AdminRepo {
	return adminRepo
}

func (r *AdminRepo) RegisterAdmin(newAdmin *Admin) (*Admin, error) {

	column := map[string]interface{}{
		"password": newAdmin.Password,
		"email":    newAdmin.Email,
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
			email,
			password
		`).
		Values(values...).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new admin register query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}
	// Execute the SQL query and get the result
	var insAdmin Admin
	err = GetReadDB().QueryRow(qry, args...).Scan(&insAdmin.Password, &insAdmin.Email)
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

	return &insAdmin, nil

}

func (r *AdminRepo) DeleteAdmin(email string) error {

	qry, args, err := GetQueryBuilder().Delete(r.Table).Where(sq.Eq{"email": email}).Where(sq.Eq{"is_superadmin": false}).ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new admin delete query",
			logger.Extra(map[string]interface{}{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return err
	}

	// Execute the DELETE query
	_, err = GetReadDB().Exec(qry, args...)
	if err != nil {
		slog.Error(
			"Failed to execute admin delete query",
			logger.Extra(map[string]interface{}{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return err
	}

	return nil
}

func (r *AdminRepo) GetFetchAdmin() ([]*SubAdmin, error) {

	qry, args, err := GetQueryBuilder().Select("email").
		From(r.Table).
		Where(sq.Eq{"is_superadmin": false}).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create Get sub-admin query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	subAdmin := []*SubAdmin{}
	err = GetReadDB().Select(&subAdmin, qry, args...)
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

	return subAdmin, nil
}

func (r *AdminRepo) ValidateAdminLogin(admin LoginReaderCredintials) (string, error) {

	qry, args, err := GetQueryBuilder().Select("is_superadmin").From(r.Table).
		Where(sq.Eq{"email": admin.Email, "password": admin.Password}).ToSql()
	if err != nil {
		slog.Error(
			"Failed to create validate login query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return "", err
	}

	// Execute the query
	var role bool
	err = GetReadDB().Get(&role, qry, args...)
	if err != nil {
		slog.Error(
			"Failed to validate login",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return "", err
	}
	fmt.Println(role)
	if role {
		return "super_admin", nil
	}
	return "admin", nil
}
