package store

import (
	"be/pkg/errors"
	"be/pkg/events"
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type UserLogsQueue interface {
	Enqueue(ctx context.Context, ev events.UserLogsEvent) error
}

type sqsService struct {
	queueURL  string
	sqsClient *sqs.Client
}

func NewUserLogsSQS(awsConfig aws.Config, queueURL string) UserLogsQueue {
	return &sqsService{sqsClient: sqs.NewFromConfig(awsConfig), queueURL: queueURL}
}

func (s *sqsService) Enqueue(ctx context.Context, ev events.UserLogsEvent) error {
	bts, err := json.Marshal(ev)
	if err != nil {
		return errors.WithStack(err)
	}

	in := &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(bts)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"route": {
				DataType:    aws.String("String"),
				StringValue: aws.String("userloggers"),
			},
		},
	}

	_, err = s.sqsClient.SendMessage(ctx, in)
	return errors.WithStack(err)
}
