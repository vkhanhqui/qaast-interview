package transport

import (
	"be/pkg/model"
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"net/url"
	"strconv"
	"time"
)

type AdminUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminListUsersInput struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
}

func (req *AdminListUsersInput) Bind(values url.Values) {
	req.Limit = 10
	if s := values.Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			req.Limit = v
		}
	}

	req.Cursor = values.Get("cursor")
}

type AdminListUsersResponse struct {
	Users      []AdminUserResponse `json:"users"`
	NextCursor string              `json:"next_cursor"`
}

func (res *AdminListUsersResponse) Bind(users []model.User) {
	if len(users) == 0 {
		res.Users = []AdminUserResponse{}
		return
	}

	res.Users = make([]AdminUserResponse, 0, len(users))
	for _, u := range users {
		res.Users = append(res.Users, AdminUserResponse{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name.String,
			CreatedAt: u.CreatedAt,
		})
	}

	res.NextCursor = res.Users[len(res.Users)-1].ID
}

type AdminUpdateUserInput struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (req *AdminUpdateUserInput) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if len(req.Email) == 0 {
		return nil
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

type AdminUpdateUserResponse struct {
	AdminUserResponse
}

func (res *AdminUpdateUserResponse) Bind(u *model.User) {
	res.ID = u.ID
	res.Email = u.Email
	res.Name = u.Name.String
	res.CreatedAt = u.CreatedAt
}

type AdminUserLogResponse struct {
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminListUserLogsInput struct {
	Limit  int    `json:"limit"`
	Cursor string `json:"cursor"`
}

func (req *AdminListUserLogsInput) Bind(values url.Values) {
	req.Limit = 10
	if s := values.Get("limit"); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			req.Limit = v
		}
	}

	req.Cursor = values.Get("cursor")
}

type AdminListUserLogsResponse struct {
	UserLogs   []AdminUserLogResponse `json:"user_logs"`
	NextCursor string                 `json:"next_cursor"`
}

func (res *AdminListUserLogsResponse) Bind(userlogs []model.UserLogs, nextCursor string) {
	if len(userlogs) == 0 {
		res.UserLogs = []AdminUserLogResponse{}
		return
	}

	res.UserLogs = make([]AdminUserLogResponse, 0, len(userlogs))
	for _, u := range userlogs {
		res.UserLogs = append(res.UserLogs, AdminUserLogResponse(u))
	}

	res.NextCursor = nextCursor
}
