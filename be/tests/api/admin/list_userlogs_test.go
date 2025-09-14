package admin

import (
	"be/tests/tester"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdminUserLogs(t *testing.T) {
	api := tester.NewAPITester()
	_, _, token := generateUser(t)

	res, err := api.Get("/admin/userlogs").
		SetHeader("Authorization", "Bearer "+token).
		Expect(t).
		Status(http.StatusOK).
		Send()
	assert.NoError(t, err)

	var logsResp adminUserLogsResp
	err = res.JSON(&logsResp)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(logsResp.UserLogs), 1)
}

type adminUserLogsResp struct {
	UserLogs []struct {
		UserID    string    `json:"user_id"`
		EventType string    `json:"event_type"`
		Details   string    `json:"details"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"user_logs"`
	NextCursor string `json:"next_cursor"`
}
