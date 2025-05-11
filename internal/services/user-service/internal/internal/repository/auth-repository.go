package repository

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"ndx/internal/services/user-service/internal/internal/dto"
	pb "ndx/pkg/api/user-service"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func (ar *AuthRepository) Register(req dto.RegisterRequest) (*pb.RegisterResponse, error) {
	queryBuilder := squirrel.Insert("users")
}
