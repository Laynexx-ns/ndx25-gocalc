package handlers

import (
	"context"
	"database/sql"
	"ndx/internal/services/user-service/internal/internal/repository"
	pb "ndx/pkg/api/user-service"
)

type AuthHandler struct {
	db   *sql.DB
	repo *repository.AuthRepository
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		db:   db,
		repo: repository.NewAuthRepository(db),
	}
}

func (ah *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

}
