package service

import (
	"api/store"
	"be/pkg/events"
	"be/pkg/model"
	"context"
	"fmt"
	"time"
)

type adminServiceWithQueue struct {
	adminSvc     AdminService
	userLogQueue store.UserLogsQueue
}

func NewAdminServiceWithQueue(adminSvc AdminService, userLogQueue store.UserLogsQueue) AdminService {
	return &adminServiceWithQueue{adminSvc: adminSvc, userLogQueue: userLogQueue}
}

func (svc *adminServiceWithQueue) ListUsers(ctx context.Context, limit int, cursor string) ([]model.User, error) {
	return svc.adminSvc.ListUsers(ctx, limit, cursor)
}

func (svc *adminServiceWithQueue) ListUserLogs(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error) {
	return svc.adminSvc.ListUserLogs(ctx, limit, cursor)
}

func (svc *adminServiceWithQueue) UpdateUser(ctx context.Context, adminID, userID, email, name string) (*model.User, error) {
	u, err := svc.adminSvc.UpdateUser(ctx, adminID, userID, email, name)
	if err != nil {
		return nil, err
	}

	err = svc.userLogQueue.Enqueue(ctx, events.UserLogsEvent{
		UserID:    u.ID,
		EventType: "admin.updateUser",
		EventTime: time.Now().UTC(),
		Details: fmt.Sprintf(
			"Admin %s updated user.\nOld values: email=%s, name=%s; New values: email=%s, name=%s",
			adminID, email, name, u.Email, u.Name.String,
		),
	})
	return u, err
}

func (svc *adminServiceWithQueue) DeleteUser(ctx context.Context, adminID, userID string) error {
	err := svc.adminSvc.DeleteUser(ctx, adminID, userID)
	if err != nil {
		return err
	}

	err = svc.userLogQueue.Enqueue(ctx, events.UserLogsEvent{
		UserID:    userID,
		EventType: "admin.deleteUser",
		EventTime: time.Now().UTC(),
		Details: fmt.Sprintf(
			"Admin %s deleted user %s",
			adminID, userID,
		),
	})
	return err
}
