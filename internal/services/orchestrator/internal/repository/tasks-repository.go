package repository

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

func (tr *TasksRepository) GetTask() *models.PrimeEvaluation {
	queryBuilder := squirrel.Select("*").
		From("prime_evaluations").
		Where(squirrel.Expr("completed_at IS NULL"))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Fatalf("can't generate query | err: %v", err)
		return nil
	}
	row := tr.db.QueryRow(query, args)

	var pe models.PrimeEvaluation
	if err = row.Scan(
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
	return &pe
}
