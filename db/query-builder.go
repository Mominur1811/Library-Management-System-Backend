package db

import (
	"github.com/Masterminds/squirrel"
)

var psql squirrel.StatementBuilderType

func InitQueryBuilder() {
	psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func GetQueryBuilder() squirrel.StatementBuilderType {
	return psql
}
