package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	examplev1 "github.com/tikayere/userservice/gen/example/v1"
	"github.com/tikayere/userservice/internal/repository"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
	logger    *zap.Logger
}

func NewUserService(repo *repository.UserRepository, jwtSecret string, logger *zap.Logger) *UserService {
	return &UserService{repo: repo, jwtSecret: jwtSecret, logger: logger}
}

func (s *UserService) RegisterUser(ctx context.Context, req *examplev1.RegisterUserRequest) (*examplev1.RegisterUserResponse, error) {
	user := &repository.User{
		ID:       fmt.Sprintf("user-%d", time.Now().UnixNano()),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Role:     req.Role,
		TenantID: req.TenantId,
	}
	if err := s.repo.CreateUser(user); err != nil {
		s.logger.Error("Failed to register user", zap.Error(err))
		return nil, err
	}
	return &examplev1.RegisterUserResponse{UserId: user.ID}, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *examplev1.LoginUserRequest) (*examplev1.LoginUserResponse, error) {
	user, err := s.repo.FindUserByEmail(req.Email)
	if err != nil {
		s.logger.Error("User not found", zap.Error(err))
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Error("Invalid password", zap.Error(err))
		return nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":   user.ID,
		"tenant_id": user.TenantID,
		"role":      user.Role,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, err
	}
	return &examplev1.LoginUserResponse{Token: tokenString}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, req *examplev1.GetUserProfileRequest) (*examplev1.GetUserProfileResponse, error) {
	user, err := s.repo.FindUserByID(req.UserId, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to get user profile", zap.Error(err))
		return nil, err
	}
	return &examplev1.GetUserProfileResponse{
		UserId:   user.ID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		TenantId: user.TenantID,
	}, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, req *examplev1.UpdateUserProfileRequest) (*examplev1.UpdateUserProfileResponse, error) {
	user, err := s.repo.FindUserByID(req.TenantId, req.TenantId)
	if err != nil {
		s.logger.Error("Failed to find user", zap.Error(err))
		return nil, err
	}
	user.Name = req.Name
	user.Email = req.Email
	if err := s.repo.UpdateUser(user); err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, err
	}
	return &examplev1.UpdateUserProfileResponse{UserId: user.ID}, nil
}

func (s *UserService) CheckPermission(ctx context.Context, req *examplev1.CheckPermissionRequest) (*examplev1.CheckPermissionResponse, error) {
	allowed, err := s.repo.CheckPermission(req.UserId, req.TenantId, req.Permission)
	if err != nil {
		s.logger.Error("Failed to check permission", zap.Error(err))
		return nil, err
	}
	return &examplev1.CheckPermissionResponse{Allowed: allowed}, nil
}
