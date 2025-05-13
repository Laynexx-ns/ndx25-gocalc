package repo

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

func (er *ExpressionRepository) GetAllExpressions() ([]models.Expressions, error) {
	queryBuilder := squirrel.Select("*").
		From("evaluations").PlaceholderFormat(squirrel.Dollar)
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
			&es.UserId,
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

func (r *ExpressionRepository) GetExpressions(userID uuid.UUID) ([]models.Expressions, error) {
	rows, err := r.db.Query(`SELECT * FROM evaluations WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Expressions
	for rows.Next() {
		var e models.Expressions
		err = rows.Scan(&e.Id, &e.UserId, &e.Status, &e.Result, &e.UserId)
		if err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func (r *ExpressionRepository) UpdateExpressionStatusAndResult(id int, status string, result float64) error {
	query := `UPDATE evaluations 
        SET status = $1, result = $2 
        WHERE id = $3`
	_, err := r.db.Exec(query, status, result, id)
	return err
}

func (r *ExpressionRepository) SaveExpression(expr models.Expressions) (int, error) {
	var id int
	err := r.db.QueryRow(`INSERT INTO evaluations (user_id, status, result, expression) VALUES ($1, $2, $3, $4) RETURNING id`, expr.UserId, expr.Status, 0, expr.Expression).Scan(&id)
	return id, err
}

func (r *ExpressionRepository) GetExpressionById(id int) (models.Expressions, error) {
	var e models.Expressions
	err := r.db.QueryRow(`SELECT id, expression, result, status FROM evaluations WHERE id = $1`, id).
		Scan(&e.Id, &e.Expression, &e.Result, &e.Status)
	return e, err
}
