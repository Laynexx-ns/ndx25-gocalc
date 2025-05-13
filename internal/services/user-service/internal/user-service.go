package internal

import (
	"context"
	"database/sql"
	"ndx/internal/services/user-service/internal/internal/handlers"
	us "ndx/pkg/api/user-service"
)

type UserService struct {
	us.UnimplementedUserServiceServer
	db          *sql.DB
	authHandler *handlers.AuthHandler
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db:          db,
		authHandler: handlers.NewAuthHandler(db),
	}
}

func (us *UserService) Register(ctx context.Context, req *us.RegisterRequest) (*us.RegisterResponse, error) {
	return us.authHandler.Register(ctx, req)
}

func (us *UserService) Login(ctx context.Context, req *us.LoginRequest) (*us.LoginResponse, error) {
	return us.authHandler.Login(ctx, req)
}
