package main

import (
	"be/pkg/config"
	"be/pkg/log"
	"be/pkg/transport/sqs"
	"context"
	"worker/repository"
	"worker/service"
	"worker/transport"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

type envConfig struct {
	Stage string `mapstructure:"STAGE"`
	Port  int    `mapstructure:"PORT"`

	AwsRegion           string `mapstructure:"AWS_REGION"`
	AWSEndpoint         string `mapstructure:"AWS_ENDPOINT"`
	AwsAccessKey        string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey        string `mapstructure:"AWS_SECRET_KEY"`
	SQSUserLogsQueueURL string `mapstructure:"SQS_USER_LOGS_QUEUE_URL"`

	DynamoTable    string `mapstructure:"DYNAMO_TABLE"`
	DynamoEndpoint string `mapstructure:"DYNAMO_ENDPOINT"`
}

func main() {
	logger := log.NewZapLogger()
	defer log.PanicRecover(logger)

	config.LoadEnvConfig()

	env := envConfig{}
	err := config.UnmarshalEnvConfig(&env)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	dynamodbAWSConfig, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(env.AwsRegion),
		awsconfig.WithBaseEndpoint(env.DynamoEndpoint),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AwsAccessKey, env.AwsSecretKey, "")))
	if err != nil {
		panic(err)
	}

	r := repository.NewLogRepo(dynamodbAWSConfig, env.DynamoTable)
	svc := service.NewLogService(r)

	userLoggersHandler := transport.NewHandler(svc)

	router := sqs.NewSQSRouter(sqs.RouteFromAttributeFn)
	router.AddHandler("userloggers", userLoggersHandler)

	sqsAWSConfig, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(env.AwsRegion),
		awsconfig.WithBaseEndpoint(env.AWSEndpoint),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AwsAccessKey, env.AwsSecretKey, "")))
	if err != nil {
		panic(err)
	}

	sqsCfg := sqs.Config{
		AWSConfig:         sqsAWSConfig,
		QueueURL:          env.SQSUserLogsQueueURL,
		MaxMessages:       10,
		VisibilityTimeout: 300,
	}
	worker := sqs.NewWorker(sqsCfg, router, logger)
	sqs.ListenForTermination(ctx, worker, logger)

	err = worker.Run(ctx)
	if err != nil {
		panic(err)
	}
}
