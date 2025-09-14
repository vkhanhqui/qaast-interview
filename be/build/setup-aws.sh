#!/bin/bash

export AWS_PAGER="" # turn off annoying pager of AWS CLI output

aws sqs create-queue \
    --queue-name user-logs-queue \
    --attributes VisibilityTimeout=30,MessageRetentionPeriod=1209600 \
    --endpoint-url "$AWS_ENDPOINT"
