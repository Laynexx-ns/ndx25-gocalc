package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/services/user-service/internal/internal/dto"
	"ndx/pkg/logger"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) Register(req dto.RegisterRequest) error {
	queryBuilder := squirrel.Insert("users").
		Columns("id", "email", "hash", "created_at").
		Values(1, req.Id).
		Values(2, req.Email).
		Values(3, req.Hash).
		Values(4, req.CreatedAt)
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if _, err = ar.db.Exec(query, args...); err != nil {
		logger.L().Logf(0, "ISE | err: %v", err)
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (ar *AuthRepository) GetUserHash(email string) (dto.AuthResponse, error) {
	query, args, err := (squirrel.Select("id, hash").
		From("users").
		Where(squirrel.Eq{"email": email})).ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return dto.AuthResponse{}, err
	}

	row := *ar.db.QueryRow(query, args...)
	var resp dto.AuthResponse
	if err = row.Scan(&resp.Id, &resp.Hash); err != nil {
		return dto.AuthResponse{}, err
	}

	return resp, nil
}
