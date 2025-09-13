package events

import "time"

type UserLogsEvent struct {
	UserID    string    `json:"userId"`
	EventType string    `json:"eventType"`
	EventTime time.Time `json:"eventTime"`
	Details   string    `json:"details"`
}
