package service

import (
	"api/store"
	"be/pkg/events"
	"be/pkg/model"
	"context"
	"fmt"
	"time"
)

type AdminService interface {
	ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error)
	ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error)
	UpdateUser(ctx context.Context, adminID, userID, email, name string) (*model.User, error)
}

type adminService struct {
	users        store.UserRepository
	userLogs     store.LogRepository
	userLogQueue store.UserLogsQueue
}

func NewAdminService(u store.UserRepository, userLogs store.LogRepository, userLogQueue store.UserLogsQueue) AdminService {
	return &adminService{users: u, userLogs: userLogs, userLogQueue: userLogQueue}
}

func (svc *adminService) ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error) {
	return svc.users.List(ctx, limit, cursor)
}

func (svc *adminService) ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error) {
	return svc.userLogs.List(ctx, limit, cursor)
}

func (svc *adminService) UpdateUser(ctx context.Context, adminID, userID, email, name string) (*model.User, error) {
	u, err := svc.users.UpdateUser(ctx, userID, email, name)
	if err != nil {
		return nil, err
	}

	if err := svc.userLogQueue.Enqueue(ctx, events.UserLogsEvent{
		UserID:    u.ID,
		EventType: "users.updateUser",
		EventTime: time.Now().UTC(),
		Details: fmt.Sprintf(
			"Admin %s updated user.\nOld values: email=%s, name=%s; New values: email=%s, name=%s",
			adminID, email, name, u.Email, u.Name.String,
		),
	}); err != nil {
		return nil, err
	}

	return u, nil
}
