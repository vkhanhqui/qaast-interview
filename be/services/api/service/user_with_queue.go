package service

import (
	"api/store"
	"be/pkg/events"
	"context"
	"fmt"
	"time"
)

type userServiceWithQueue struct {
	svc          UserService
	userLogQueue store.UserLogsQueue
}

func NewUserServiceWithQueue(svc UserService, userLogQueue store.UserLogsQueue) UserService {
	return &userServiceWithQueue{svc: svc, userLogQueue: userLogQueue}
}

func (s *userServiceWithQueue) SignUp(ctx context.Context, email, password string) (string, error) {
	id, err := s.svc.SignUp(ctx, email, password)
	if err != nil {
		return "", err
	}

	if err := s.userLogQueue.Enqueue(ctx, events.UserLogsEvent{
		UserID:    id,
		EventType: "users.signUp",
		EventTime: time.Now().UTC(),
		Details:   fmt.Sprintf("New user: id=%s email=%s", id, email),
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (s *userServiceWithQueue) SignIn(ctx context.Context, email, password string) (string, string, error) {
	ss, id, err := s.svc.SignIn(ctx, email, password)
	if err != nil {
		return "", "", err
	}

	if err := s.userLogQueue.Enqueue(ctx, events.UserLogsEvent{
		UserID:    id,
		EventType: "users.signIn",
		EventTime: time.Now().UTC(),
		Details:   fmt.Sprintf("User signed-in: email=%s", email),
	}); err != nil {
		return "", "", err
	}

	return ss, id, nil
}
