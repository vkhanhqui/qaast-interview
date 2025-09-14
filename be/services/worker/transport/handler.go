package transport

import (
	"be/pkg/errors"
	"be/pkg/events"
	"be/pkg/transport/sqs"
	"context"
	"encoding/json"
	"worker/service"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func NewHandler(svc service.LogService) sqs.HandlerFunc {
	h := &handler{svc}
	return h.HandlerFunc
}

type handler struct {
	svc service.LogService
}

func (h *handler) HandlerFunc(ctx context.Context, msg types.Message) error {
	ev := events.UserLogsEvent{}
	if err := json.Unmarshal([]byte(aws.ToString(msg.Body)), &ev); err != nil {
		return errors.WithStack(err)
	}

	return h.svc.Write(ctx, ev)
}
