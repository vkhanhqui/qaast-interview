package transport

import (
	"be/pkg/model"
	"time"
)

type SignUpInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	UserID string `json:"user_id"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

type UpdateUserInput struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type UpdateUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

func (res *UpdateUserResponse) Bind(u *model.User) {
	res.ID = u.ID
	res.Email = u.Email
	res.Name = u.Name.String
	res.CreatedAt = u.CreatedAt
}

type ErrorResponse struct {
	Error string `json:"error"`
}
