package repo

import (
	"database/sql"
	"github.com/google/uuid"
	"ndx/internal/models"
)

type ExpressionRepository struct {
	db *sql.DB
}

func NewExpressionRepository(db *sql.DB) *ExpressionRepository {
	return &ExpressionRepository{
		db: db,
	}
}

//func (er *ExpressionRepository) GetExpressions(userId uuid.UUID) ([]models.Expressions, error) {
//	queryBuilder := squirrel.Select("*").
//		From("evaluations").
//		Where(squirrel.Eq{"user_id": userId})
//	query, args, err := queryBuilder.ToSql()
//	if err != nil {
//		logger.L().Logf(0, "can't build query | err: %v", err)
//	}
//
//	rows, err := er.db.Query(query, args...)
//	if err != nil {
//		logger.L().Logf(0, "ISE | err: %v", err)
//		return []models.Expressions{}, nil
//	}
//
//	var expressions []models.Expressions
//	for rows.Next() {
//		var es models.Expressions
//		if err = rows.Scan(
//			&es.Id,
//			&es.UserId,
//			&es.Status,
//			&es.Result,
//			&es.Expression,
//		); err != nil {
//			logger.L().Logf(0, "ISE | err: %v", err)
//		}
//		expressions = append(expressions, es)
//	}
//
//	return expressions, nil
//}

func (r *ExpressionRepository) GetExpressions(userID uuid.UUID) ([]models.Expressions, error) {
	rows, err := r.db.Query(`SELECT id, expression, result, status, user_id FROM evaluations WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Expressions
	for rows.Next() {
		var e models.Expressions
		err := rows.Scan(&e.Id, &e.Expression, &e.Result, &e.Status, &e.UserId)
		if err != nil {
			return nil, err
		}
		results = append(results, e)
	}
	return results, nil
}

func (r *ExpressionRepository) SaveExpression(expr models.Expressions) (int, error) {
	var id int
	err := r.db.QueryRow(`INSERT INTO evaluations (expression, status, user_id) VALUES ($1, $2, $3) RETURNING id`, expr.Expression, expr.Status, expr.UserId).Scan(&id)
	return id, err
}

func (r *ExpressionRepository) GetExpressionById(id int) (models.Expressions, error) {
	var e models.Expressions
	err := r.db.QueryRow(`SELECT id, expression, result, status FROM evaluations WHERE id = $1`, id).
		Scan(&e.Id, &e.Expression, &e.Result, &e.Status)
	return e, err
}
