package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/services/user-service/internal/internal/dto"
	"ndx/pkg/logger"
)

var psql = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) Register(req dto.RegisterRequest) error {
	exists, err := ar.checkIfUserExists(req.Email)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	if exists {
		return status.Error(codes.AlreadyExists, "exists")
	}

	queryBuilder := psql.Insert("users").
		Columns("id", "email", "hash", "created_at").
		Values(req.Id, req.Email, req.Hash, req.CreatedAt)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		logger.L().Logf(0, "can't build query | err: %v", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}

	logger.L().Logf(0, "query: %s, args: %v", query, args)

	if _, err = ar.db.Exec(query, args...); err != nil {
		logger.L().Logf(0, "ISE | err: %v", err)
		return status.Error(codes.Internal, err.Error())
	}
	return nil
}

func (ar *AuthRepository) checkIfUserExists(email string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 from users WHERE email = $1)`

	var exists bool
	if err := ar.db.QueryRow(query, email).Scan(&exists); err != nil {
		logger.L().Logf(0, "ISE | err: %v", err)
		return false, err
	}

	if exists {
		return true, nil
	}
	return false, nil

}

func (ar *AuthRepository) GetUserHash(email string) (dto.AuthResponse, error) {
	query, args, err := (psql.Select("id, hash").
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
