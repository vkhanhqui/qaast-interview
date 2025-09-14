package admin

import (
	"be/tests/tester"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdminListUsers(t *testing.T) {
	api := tester.NewAPITester()
	_, _, token := generateUser(t)

	res, err := api.Get("/admin/users").
		AddQuery("limit", "1000").
		SetHeader("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		Send()
	assert.NoError(t, err)

	var adminResp adminUsersResp
	err = res.JSON(&adminResp)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(adminResp.Users), 1)
}

func generateUser(t *testing.T) (userID, email, token string) {
	api := tester.NewAPITester()

	email = fmt.Sprintf("test+%d@example.com", time.Now().UnixNano())
	password := "test"
	res, err := api.Post("/users/signup").
		SetHeader("Content-Type", "application/json").
		BodyString(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)).
		Expect(t).
		Status(http.StatusCreated).
		Send()
	assert.NoError(t, err)

	var signUpResp signUpResponse
	err = res.JSON(&signUpResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, signUpResp.UserID)

	res, err = api.Post("/users/signin").
		SetHeader("Content-Type", "application/json").
		BodyString(fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)).
		Expect(t).
		Status(http.StatusOK).
		Send()
	assert.NoError(t, err)

	var tokenResp signInResponse
	err = res.JSON(&tokenResp)
	assert.NoError(t, err)
	assert.NotEmpty(t, tokenResp.Token)

	return signUpResp.UserID, email, tokenResp.Token
}

type signUpResponse struct {
	UserID string `json:"user_id"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type adminUsersResp struct {
	Users []struct {
		ID        string    `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"users"`
	NextCursor string `json:"next_cursor"`
}
