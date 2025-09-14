package transport

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
)

type SignUpInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *SignUpInput) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

type SignUpResponse struct {
	UserID string `json:"user_id"`
}

type SignInInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *SignInInput) Bind(r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return errors.New("invalid email format")
	}

	return nil
}

type SignInResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
