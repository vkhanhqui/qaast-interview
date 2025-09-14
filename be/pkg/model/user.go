package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string
	Email     string
	Password  string
	Name      sql.NullString
	CreatedAt time.Time
}
