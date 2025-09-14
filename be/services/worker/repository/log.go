package repository

import (
	"be/pkg/errors"
	"be/pkg/model"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbtypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type LogRepository interface {
	Write(ctx context.Context, l model.UserLogs) error
}

type logRepo struct {
	client *dynamodb.Client
	table  string
}

func NewLogRepo(cfg aws.Config, table string) LogRepository {
	return &logRepo{
		client: dynamodb.NewFromConfig(cfg),
		table:  table,
	}
}

func (r *logRepo) Write(ctx context.Context, l model.UserLogs) error {
	id, err := uuid.NewV7()
	if err != nil {
		return errors.WithStack(err)
	}

	item := map[string]ddbtypes.AttributeValue{
		"PK":         &ddbtypes.AttributeValueMemberS{Value: "logs"},
		"SK":         &ddbtypes.AttributeValueMemberS{Value: fmt.Sprintf("%s#%s", l.UserID, id)},
		"user_id":    &ddbtypes.AttributeValueMemberS{Value: l.UserID},
		"event_type": &ddbtypes.AttributeValueMemberS{Value: l.EventType},
		"details":    &ddbtypes.AttributeValueMemberS{Value: l.Details},
		"created_at": &ddbtypes.AttributeValueMemberS{Value: l.CreatedAt.Format(time.RFC3339Nano)},
	}
	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.table,
		Item:      item,
	})
	return errors.WithStack(err)
}
