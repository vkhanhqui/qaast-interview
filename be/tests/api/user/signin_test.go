package api

import (
	"be/tests/tester"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignInAPI(t *testing.T) {
	api := tester.NewAPITester()

	email := fmt.Sprintf("test+%d@example.com", time.Now().UnixNano())
	password := "test"
	err := api.Post("/users/signup").
		SetHeader("Content-Type", "application/json").
		BodyString(fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)).
		Expect(t).
		Status(http.StatusCreated).
		Done()
	assert.NoError(t, err)

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "invalid email format",
			body:           `{"email":"test111112222","password":"test"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "valid credentials",
			body:           fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "wrong password",
			body:           fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, "wrong password"),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = api.Post("/users/signin").
				SetHeader("Content-Type", "application/json").
				BodyString(tt.body).
				Expect(t).
				Status(tt.expectedStatus).
				Done()
			assert.NoError(t, err)
		})
	}
}
