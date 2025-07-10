package handler

import (
	"context"

	examplev1 "github.com/tikayere/userservice/gen/example/v1"
	"github.com/tikayere/userservice/internal/service"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc    *service.UserService
	logger *zap.Logger
	examplev1.UnimplementedUserServiceServer
}

func NewUserHandler(svc *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *examplev1.RegisterUserRequest) (*examplev1.RegisterUserResponse, error) {
	return h.svc.RegisterUser(ctx, req)
}

func (h *UserHandler) LoginUser(ctx context.Context, req *examplev1.LoginUserRequest) (*examplev1.LoginUserResponse, error) {
	return h.svc.LoginUser(ctx, req)
}

func (h *UserHandler) GetUserProfile(ctx context.Context, req *examplev1.GetUserProfileRequest) (*examplev1.GetUserProfileResponse, error) {
	return h.svc.GetUserProfile(ctx, req)
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *examplev1.UpdateUserProfileRequest) (*examplev1.UpdateUserProfileResponse, error) {
	return h.svc.UpdateUserProfile(ctx, req)
}

func (h *UserHandler) CheckPermission(ctx context.Context, req *examplev1.CheckPermissionRequest) (*examplev1.CheckPermissionResponse, error) {
	return h.svc.CheckPermission(ctx, req)
}
