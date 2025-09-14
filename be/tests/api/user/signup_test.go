package api

import (
	"be/tests/tester"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSignupAPI(t *testing.T) {
	api := tester.NewAPITester()

	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "invalid email",
			body:           `{"email": "test", "password": "test"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "valid email",
			body:           fmt.Sprintf(`{"email": "test+%d@example.com", "password": "test"}`, time.Now().UnixNano()),
			expectedStatus: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := api.Post("/users/signup").
				SetHeader("Content-Type", "application/json").
				BodyString(tt.body).
				Expect(t).
				Status(tt.expectedStatus).
				Done()
			assert.NoError(t, err)
		})
	}
}
