package service

import (
	"be/pkg/errors"
	"be/pkg/events"
	"be/pkg/model"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type userServiceSQS struct {
	svc       UserService
	sqsClient *sqs.Client
	queueURL  string
}

func NewUserServiceSQS(svc UserService, awsConfig aws.Config, queueURL string) UserService {
	return &userServiceSQS{svc: svc, sqsClient: sqs.NewFromConfig(awsConfig), queueURL: queueURL}
}

func (s *userServiceSQS) SignUp(ctx context.Context, email, password string) (string, error) {
	id, err := s.svc.SignUp(ctx, email, password)
	if err != nil {
		return "", err
	}

	if err := s.enqueue(ctx, events.UserLogsEvent{
		UserID:    id,
		EventType: "users.signUp",
		EventTime: time.Now().UTC(),
		Details:   fmt.Sprintf("New user: id=%s email=%s", id, email),
	}); err != nil {
		return "", err
	}

	return id, nil
}

func (s *userServiceSQS) SignIn(ctx context.Context, email, password string) (string, string, error) {
	ss, id, err := s.svc.SignIn(ctx, email, password)
	if err != nil {
		return "", "", err
	}

	if err := s.enqueue(ctx, events.UserLogsEvent{
		UserID:    id,
		EventType: "users.signIn",
		EventTime: time.Now().UTC(),
		Details:   fmt.Sprintf("User signed-in: email=%s", email),
	}); err != nil {
		return "", "", err
	}

	return ss, id, nil
}

func (s *userServiceSQS) UpdateUser(ctx context.Context, userID, email, name string) (*model.User, error) {
	u, err := s.svc.UpdateUser(ctx, userID, email, name)
	if err != nil {
		return nil, err
	}

	if err := s.enqueue(ctx, events.UserLogsEvent{
		UserID:    u.ID,
		EventType: "users.updateUser",
		EventTime: time.Now().UTC(),
		Details: fmt.Sprintf(
			"Old values: email=%s, name=%s; New values: email=%s, name=%s",
			email, name, u.Email, u.Name.String,
		),
	}); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *userServiceSQS) enqueue(ctx context.Context, ev events.UserLogsEvent) error {
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
