package store

import (
	"be/pkg/errors"
	"be/pkg/model"
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type LogRepository interface {
	List(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error)
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

func (r *logRepo) List(ctx context.Context, limit int, cursor string) ([]model.UserLogs, string, error) {
	params := &dynamodb.QueryInput{
		TableName:              &r.table,
		KeyConditionExpression: aws.String("PK = :pk"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk": &types.AttributeValueMemberS{Value: "logs"},
		},
		Limit:            aws.Int32(int32(limit)),
		ScanIndexForward: aws.Bool(false),
	}

	if len(cursor) > 0 {
		params.ExclusiveStartKey = map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: "logs"},
			"SK": &types.AttributeValueMemberS{
				Value: cursor,
			},
		}
	}

	out, err := r.client.Query(ctx, params)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	logs := make([]model.UserLogs, 0, len(out.Items))
	for _, it := range out.Items {
		createdAt, err := time.Parse(time.RFC3339Nano, it["created_at"].(*types.AttributeValueMemberS).Value)
		if err != nil {
			return nil, "", errors.WithStack(err)
		}

		logs = append(logs, model.UserLogs{
			UserID:    it["user_id"].(*types.AttributeValueMemberS).Value,
			EventType: it["event_type"].(*types.AttributeValueMemberS).Value,
			Details:   it["details"].(*types.AttributeValueMemberS).Value,
			CreatedAt: createdAt,
		})
	}

	var nextCursor string
	if skAttr, ok := out.LastEvaluatedKey["SK"].(*types.AttributeValueMemberS); ok {
		nextCursor = skAttr.Value
	}

	return logs, nextCursor, nil
}
