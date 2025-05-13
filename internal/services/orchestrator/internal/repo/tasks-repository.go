package repo

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"ndx/internal/models"
	"ndx/pkg/logger"
)

type TasksRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TasksRepository {
	return &TasksRepository{
		db: db,
	}
}

func (tr *TasksRepository) GetAllTasks() []models.PrimeEvaluation {
	queryBuilder := squirrel.Select("*").
		From("prime_evaluations")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Fatalf("can't generate query | err: %v", err)
		return nil
	}
	rows, err := tr.db.Query(query, args)
	if err != nil {
		logger.L().Logf(0, "ISE | err: %v", err)
	}
	var evals []models.PrimeEvaluation
	for rows.Next() {
		var pe models.PrimeEvaluation
		if err = rows.Scan(
			&pe.ParentID,
			&pe.Id,
			&pe.Arg1,
			&pe.Arg2,
			&pe.Arg2,
			&pe.Operation,
			&pe.OperationTime,
			&pe.Result,
			&pe.Error,
			&pe.CompletedAt,
			&pe.UserId,
		); err != nil {
			logger.L().Logf(0, "not found or error | err: %v", err)
			return nil
		}
		evals = append(evals, pe)
	}

	return evals
}

func (tr *TasksRepository) GetPendingTasks() ([]models.Expressions, error) {
	queryBuilder := squirrel.Select("id", "expression", "user_id", "status").
		From("evaluations").
		Where(squirrel.Eq{"status": "pending"}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return nil, err
	}

	rows, err := tr.db.Query(query, args...)
	if err != nil {
		logger.L().Logf(0, "can complete query | err: %v", err)
		return nil, err
	}
	defer rows.Close()

	var expressions []models.Expressions
	for rows.Next() {
		var e models.Expressions
		if err = rows.Scan(&e.Id, &e.Expression, &e.UserId, &e.Status); err != nil {
			logger.L().Logf(0, "can't scan to model | err: %v", err)
			return nil, err
		}
		expressions = append(expressions, e)
	}

	return expressions, nil
}

func (tr *TasksRepository) UpdateExpressionResult(id int, status string, result float64) error {
	queryBuilder := squirrel.Update("evaluations").
		Set("status", status).
		Set("result", result).
		Where(squirrel.Eq{"id": id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tr.db.Exec(query, args...)
	return err
}

func (tr *TasksRepository) UpdateExpressionStatus(id int, status string) error {
	queryBuilder := squirrel.Update("evaluations").
		Set("status", status).
		Where(squirrel.Eq{"id": id})

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tr.db.Exec(query, args...)
	return err
}

func (tr *TasksRepository) SavePrimeEvaluation(pe models.PrimeEvaluation) error {
	queryBuilder := squirrel.Insert("prime_evaluations").
		Columns("parent_id", "arg1", "arg2", "operation", "operation_time", "result").
		Values(pe.ParentID, pe.Arg1, pe.Arg2, pe.Operation, pe.OperationTime, 0).PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return err
	}

	n, err := tr.db.Exec(query, args...)
	logger.L().Logf(0, "n: %v | err: %v", n, err)

	return err
}

func (tr *TasksRepository) GetPrimeEvaluationByParentID(parentID int) ([]models.PrimeEvaluation, error) {
	queryBuilder := squirrel.Select("id", "parent_id", "arg1", "arg2", "operation", "operation_time", "result", "error", "completed_at").
		From("prime_evaluations").
		Where(squirrel.Eq{"parent_id": parentID}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return nil, err
	}

	rows, err := tr.db.Query(query, args...)
	if err != nil {
		logger.L().Logf(0, "can't complete query | err: %v", err)
		return nil, err
	}
	defer rows.Close()

	logger.L().Logf(0, "rows: %v", rows)

	var results []models.PrimeEvaluation
	for rows.Next() {
		var pe models.PrimeEvaluation
		if err = rows.Scan(
			&pe.Id,
			&pe.ParentID,
			&pe.Arg1,
			&pe.Arg2,
			&pe.Operation,
			&pe.OperationTime,
			&pe.Result,
			&pe.Error,
			&pe.CompletedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, pe)
	}

	return results, nil
}
