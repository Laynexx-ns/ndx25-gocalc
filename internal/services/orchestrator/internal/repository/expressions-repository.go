package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"ndx/internal/models"
	"ndx/pkg/logger"
)

type ExpressionRepository struct {
	db *sql.DB
}

func NewExpressionRepository(db *sql.DB) *ExpressionRepository {
	return &ExpressionRepository{
		db: db,
	}
}

func (er *ExpressionRepository) GetExpressions(userId uuid.UUID) ([]models.Expressions, error) {
	queryBuilder := squirrel.Select("*").
		From("evaluations").
		Where(squirrel.Eq{"user_id": userId})
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
	}

	rows, err := er.db.Query(query, args...)
	if err != nil {
		logger.L().Logf(0, "ISE | err: %v", err)
		return []models.Expressions{}, nil
	}

	var expressions []models.Expressions
	for rows.Next() {
		var es models.Expressions
		if err = rows.Scan(
			&es.Id,
			&es.Status,
			&es.Result,
			&es.Expression,
		); err != nil {
			logger.L().Logf(0, "ISE | err: %v", err)
		}
		expressions = append(expressions, es)
	}

	return expressions, nil
}
