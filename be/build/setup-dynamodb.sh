#!/bin/bash

export AWS_PAGER="" # turn off annoying pager of AWS CLI output

aws dynamodb create-table \
  --table-name "user_activity_logs" \
  --attribute-definitions \
      AttributeName=PK,AttributeType=S \
      AttributeName=SK,AttributeType=S \
  --key-schema \
      AttributeName=PK,KeyType=HASH \
      AttributeName=SK,KeyType=RANGE \
  --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 \
  --endpoint-url "$AWS_ENDPOINT"
