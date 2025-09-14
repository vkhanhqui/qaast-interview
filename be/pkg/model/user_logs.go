package model

import (
	"time"
)

type UserLogs struct {
	UserID    string
	EventType string
	Details   string
	CreatedAt time.Time
}
