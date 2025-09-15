package service

import (
	"api/store"
	"be/pkg/errors"
	"be/pkg/model"
	"context"
)

type AdminService interface {
	ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error)
	ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error)
	UpdateUser(ctx context.Context, adminID, userID, email, name string) (*model.User, error)
	DeleteUser(ctx context.Context, adminID, userID string) error
}

type adminService struct {
	users    store.UserRepository
	userLogs store.LogRepository
}

func NewAdminService(u store.UserRepository, userLogs store.LogRepository) AdminService {
	return &adminService{users: u, userLogs: userLogs}
}

func (svc *adminService) ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error) {
	return svc.users.List(ctx, limit, cursor)
}

func (svc *adminService) ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error) {
	return svc.userLogs.List(ctx, limit, cursor)
}

func (svc *adminService) UpdateUser(ctx context.Context, adminID, userID, email, name string) (*model.User, error) {
	return svc.users.UpdateUser(ctx, userID, email, name)
}

func (svc *adminService) DeleteUser(ctx context.Context, adminID, userID string) error {
	if adminID == userID {
		return errors.WithInvalid(errors.New("Could not delete yourself"), "")
	}
	return svc.users.DeleteUser(ctx, userID)
}
