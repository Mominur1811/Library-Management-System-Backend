package db

import (
	"librarymanagement/logger"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
)

type Category struct {
	Name string `db:"category" validate:"required"  json:"category"`
}

type CategoryRepo struct {
	Table string
}

var categoryRepo *CategoryRepo

func InitCategoryRepo() {
	categoryRepo = &CategoryRepo{Table: "category"}
}

func GetCategoryRepo() *CategoryRepo {
	return categoryRepo
}

func (r *CategoryRepo) InsertCategory(newCategory *Category) (*Category, error) {

	column := map[string]interface{}{
		"category": newCategory.Name,
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
			category
		`).
		Values(values...).
		ToSql()
	if err != nil {
		slog.Error(
			"Failed to create new category query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	var insertedCategory Category
	err = GetReadDB().QueryRow(qry, args...).Scan(&insertedCategory.Name)
	if err != nil {
		slog.Error(
			"Failed to execute insert category",
			logger.Extra(map[string]interface{}{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	return &insertedCategory, nil
}

func (r *CategoryRepo) GetCategory() ([]*Category, error) {

	qry, args, err := GetQueryBuilder().Select("*").
		From(r.Table).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create get category query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	categories := []*Category{}
	err = GetReadDB().Select(&categories, qry, args...)
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

	return categories, nil
}

func (r *CategoryRepo) FindCategory(checkCategoryValue string) (*int, error) {

	qry, args, err := GetQueryBuilder().Select("Count(*)").
		From(r.Table).Where(sq.Eq{"category": checkCategoryValue}).
		ToSql()

	if err != nil {
		slog.Error(
			"Failed to create count of category query",
			logger.Extra(map[string]any{
				"error": err.Error(),
				"query": qry,
				"args":  args,
			}),
		)
		return nil, err
	}

	var count int
	err = GetReadDB().Get(&count, qry, args...)
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

	return &count, nil
}
