package service

import (
	"be/pkg/events"
	"be/pkg/model"
	"context"
	"worker/store"
)

type LogService interface {
	Write(ctx context.Context, ev events.UserLogsEvent) error
}

type logService struct {
	logRepo store.LogRepository
}

func NewLogService(r store.LogRepository) LogService {
	return &logService{logRepo: r}
}

func (s *logService) Write(ctx context.Context, ev events.UserLogsEvent) error {
	return s.logRepo.Write(ctx, model.UserLogs{
		UserID:    ev.UserID,
		EventType: ev.EventType,
		Details:   ev.Details,
		CreatedAt: ev.EventTime,
	})
}
