package service

import (
	repo "api/repository"
	"be/pkg/model"
	"context"
)

type AdminService interface {
	ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error)
	ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error)
}

type adminService struct {
	users    repo.UserRepository
	userLogs repo.LogRepository
}

func NewAdminService(u repo.UserRepository, userLogs repo.LogRepository) AdminService {
	return &adminService{users: u, userLogs: userLogs}
}

func (svc *adminService) ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error) {
	return svc.users.List(ctx, limit, cursor)
}

func (svc *adminService) ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error) {
	return svc.userLogs.List(ctx, limit, cursor)
}
