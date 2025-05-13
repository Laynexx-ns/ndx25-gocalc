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

//func (tr *TasksRepository) GetAllTasks() []models.Expressions {
//	queryBuilder := squirrel.Select("*").
//		From("prime_evaluations")
//
//	query, args, err := queryBuilder.ToSql()
//	if err != nil {
//		logger.L().Fatalf("can't generate query | err: %v", err)
//		return nil
//	}
//	rows, err := tr.db.Query(query, args)
//	if err != nil {
//		logger.L().Logf(0, "ISE | err: %v", err)
//	}
//	var evals []models.PrimeEvaluation
//	for rows.Next() {
//		var pe models.PrimeEvaluation
//		if err = rows.Scan(
//			&pe.ParentID,
//			&pe.Id,
//			&pe.Arg1,
//			&pe.Arg2,
//			&pe.Arg2,
//			&pe.Operation,
//			&pe.OperationTime,
//			&pe.Result,
//			&pe.Error,
//			&pe.CompletedAt,
//			&pe.UserId,
//		); err != nil {
//			logger.L().Logf(0, "not found or error | err: %v", err)
//			return nil
//		}
//		evals = append(evals, pe)
//	}
//
//	return evals
//}

func (tr *TasksRepository) GetPendingTasks() (models.PrimeEvaluation, error) {
	queryBuilder := squirrel.Select("*").
		From("prime_evaluations").
		Where("completed_at is NULL").PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return models.PrimeEvaluation{}, err
	}

	row := tr.db.QueryRow(query, args...)

	var expression models.PrimeEvaluation
	if err = row.Scan(
		&expression.Id,
		&expression.ParentID,
		&expression.Arg1,
		&expression.Arg2,
		&expression.Operation,
		&expression.OperationTime,
		&expression.Result,
		&expression.Error,
		&expression.CompletedAt,
	); err != nil {
		logger.L().Logf(0, "ISE | ERR: %v", err)
		return models.PrimeEvaluation{}, err
	}
	return expression, nil
}

func (tr *TasksRepository) UpdateExpressionResult(id int, status string, result float64) error {
	queryBuilder := squirrel.Update("evaluations").
		Set("status", status).
		Set("result", result).
		Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar)

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

func (tr *TasksRepository) SavePrimeEvaluation(pe models.PrimeEvaluation) (int, error) {
	queryBuilder := squirrel.Insert("prime_evaluations").
		Columns("parent_id", "arg1", "arg2", "operation", "operation_time", "result").
		Values(pe.ParentID, pe.Arg1, pe.Arg2, pe.Operation, pe.OperationTime, 0).
		PlaceholderFormat(squirrel.Dollar).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return 0, err
	}

	var id int
	err = tr.db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		logger.L().Logf(0, "can't insert evaluation | err: %v", err)
		return 0, err
	}

	logger.L().Logf(0, "inserted evaluation with id: %d", id)
	return id, nil
}

func (tr *TasksRepository) GetPrimeEvaluationByParentID(parentID int) ([]models.PrimeEvaluation, error) {
	queryBuilder := squirrel.Select(
		"id",
		"parent_id",
		"arg1",
		"arg2",
		"operation",
		"operation_time",
		"result",
		"error",
		"completed_at",
	).
		From("prime_evaluations").
		Where(squirrel.Eq{"parent_id": parentID}).
		PlaceholderFormat(squirrel.Dollar)

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

	//logger.L().Logf(0, "rows: %v", rows)

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

func (tr *TasksRepository) GetPrimeEvaluationByID(id int) (models.PrimeEvaluation, error) {
	queryBuilder := squirrel.Select("id", "parent_id", "arg1", "arg2", "operation", "operation_time", "result", "error", "completed_at").
		From("prime_evaluations").
		Where(squirrel.Eq{"id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return models.PrimeEvaluation{}, err
	}

	row := tr.db.QueryRow(query, args...)

	var res models.PrimeEvaluation
	if err = row.Scan(
		&res.Id,
		&res.ParentID,
		&res.Arg1,
		&res.Arg2,
		&res.Operation,
		&res.OperationTime,
		&res.Result,
		&res.Error,
		&res.CompletedAt,
	); err != nil {
		logger.L().Logf(0, "can't parse data from db | err: %v", err)
	}

	return res, nil
}
