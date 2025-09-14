package service

import (
	"api/store"
	"be/pkg/errors"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	SignUp(ctx context.Context, email, password string) (string, error)
	SignIn(ctx context.Context, email, password string) (string, string, error)
}

type userService struct {
	users  store.UserRepository
	jwtKey []byte
}

func NewUserService(u store.UserRepository, secret string) UserService {
	return &userService{users: u, jwtKey: []byte(secret)}
}

func (s *userService) SignUp(ctx context.Context, email, password string) (string, error) {
	_, err := s.users.FindByEmail(ctx, email)
	if !errors.IsNotFound(err) {
		if err == nil {
			err = errors.WithInvalid(errors.New("Email existed"), "")
		}
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return s.users.Create(ctx, email, string(hash))
}

func (s *userService) SignIn(ctx context.Context, email, password string) (string, string, error) {
	u, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": u.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	ss, err := token.SignedString(s.jwtKey)
	return ss, u.ID, err
}
