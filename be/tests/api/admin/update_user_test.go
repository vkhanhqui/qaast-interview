package admin

import (
	"be/tests/tester"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdminUpdateUser(t *testing.T) {
	api := tester.NewAPITester()
	_, _, adminToken := generateUser(t)

	type testCase struct {
		testCase string
		email    string
		name     string
		status   int
	}

	cases := []testCase{
		{
			testCase: "valid update name and email",
			email:    fmt.Sprintf("newemail+test+%d@example.com", time.Now().UnixNano()),
			name:     "new-name",
			status:   http.StatusOK,
		},
		{
			testCase: "missing name allowed",
			email:    fmt.Sprintf("newemail+test+%d@example.com", time.Now().UnixNano()),
			status:   http.StatusOK,
		},
		{
			testCase: "invalid email rejected",
			email:    "invalid",
			status:   http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			userID, _, _ := generateUser(t)
			updateBody := fmt.Sprintf(
				`{"id":"%s","email":"%s","name":"%s"}`,
				userID,
				tc.email,
				tc.name,
			)

			res, err := api.Put("/admin/users").
				SetHeader("Authorization", "Bearer "+adminToken).
				SetHeader("Content-Type", "application/json").
				BodyString(updateBody).
				Expect(t).
				Status(tc.status).
				Send()
			assert.NoError(t, err)

			if tc.status != http.StatusOK {
				return
			}

			var adminUpdateUserRes adminUpdateUserResponse
			err = res.JSON(&adminUpdateUserRes)
			assert.NoError(t, err)

			assert.Equal(t, adminUpdateUserRes.ID, userID)
			assert.Equal(t, adminUpdateUserRes.Email, tc.email)
			assert.Equal(t, adminUpdateUserRes.Name, tc.name)
		})
	}
}

type adminUpdateUserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
