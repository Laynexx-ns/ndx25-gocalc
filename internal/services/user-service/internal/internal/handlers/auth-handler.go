package handlers

import (
	"context"
	"database/sql"
	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"ndx/internal/services/user-service/internal/internal/dto"
	"ndx/pkg/utils"
	"time"

	"ndx/internal/services/user-service/internal/internal/repository"
	pb "ndx/pkg/api/user-service"
	"ndx/pkg/logger"
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
	if req.Password != req.PasswordConfirm {
		return nil, status.Error(codes.InvalidArgument, "passwords dont match")
	}
	hash, err := argon2id.CreateHash(req.Password, &argon2id.Params{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	})
	if err != nil {
		logger.L().Logf(0, "can't create hash of the password | err: %v", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	id := uuid.New()
	if err = ah.repo.Register(dto.RegisterRequest{
		Id:        id,
		Email:     req.Email,
		Hash:      hash,
		CreatedAt: time.Now(),
	}); err != nil {
		logger.L().Logf(0, "ISE | err: %v", err)
		return nil, status.Error(codes.Internal, err.Error())

	}

	token := utils.NewToken(req.Email)
	if token == "" {
		logger.L().Log(0, "can't create jwt token")
		return nil, status.Error(codes.Internal, "internal")
	}

	return &pb.RegisterResponse{
		Uuid:  id.String(),
		Token: token,
	}, nil

}

func (ah *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {

	resp, err := ah.repo.GetUserHash(req.Email)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	ok, err := argon2id.ComparePasswordAndHash(req.Password, resp.Hash)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	if ok {
		return &pb.LoginResponse{
			Uuid:  resp.Id.String(),
			Token: resp.Hash,
		}, nil
	}
}
