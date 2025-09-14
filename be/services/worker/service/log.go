package service

import (
	"be/pkg/events"
	"be/pkg/model"
	"context"
	repo "worker/repository"
)

type LogService interface {
	Write(ctx context.Context, ev events.UserLogsEvent) error
}

type logService struct {
	repo repo.LogRepository
}

func NewLogService(r repo.LogRepository) LogService {
	return &logService{repo: r}
}

func (s *logService) Write(ctx context.Context, ev events.UserLogsEvent) error {
	return s.repo.Write(ctx, model.UserLogs{
		UserID:    ev.UserID,
		EventType: ev.EventType,
		Details:   ev.Details,
		CreatedAt: ev.EventTime,
	})
}
