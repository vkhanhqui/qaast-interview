package main

import (
	repo "api/repository"
	"api/service"
	"api/transport"
	"be/pkg/config"
	"context"
	"fmt"
	"log"
	"net/http"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

type envConfig struct {
	Stage        string `mapstructure:"STAGE"`
	Port         int    `mapstructure:"PORT"`
	PgApiConnURI string `mapstructure:"PG_API_CONN_URI"`
	JwtSecret    string `mapstructure:"JWT_SECRET"`

	AwsRegion           string `mapstructure:"AWS_REGION"`
	AWSEndpoint         string `mapstructure:"AWS_ENDPOINT"`
	AwsAccessKey        string `mapstructure:"AWS_ACCESS_KEY"`
	AwsSecretKey        string `mapstructure:"AWS_SECRET_KEY"`
	SQSUserLogsQueueURL string `mapstructure:"SQS_USER_LOGS_QUEUE_URL"`

	DynamoTable    string `mapstructure:"DYNAMO_TABLE"`
	DynamoEndpoint string `mapstructure:"DYNAMO_ENDPOINT"`
}

func main() {
	config.LoadEnvConfig()

	env := envConfig{}
	err := config.UnmarshalEnvConfig(&env)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pgPool, err := pgxpool.New(ctx, env.PgApiConnURI)
	if err != nil {
		panic(err)
	}
	defer pgPool.Close()

	sqsAWSConfig, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(env.AwsRegion),
		awsconfig.WithBaseEndpoint(env.AWSEndpoint),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AwsAccessKey, env.AwsSecretKey, "")))
	if err != nil {
		panic(err)
	}

	dynamodbAWSConfig, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(env.AwsRegion),
		awsconfig.WithBaseEndpoint(env.DynamoEndpoint),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AwsAccessKey, env.AwsSecretKey, "")))
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://frontend:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders: []string{"Link"},
	}))

	userRepo := repo.NewUserRepo(pgPool)
	userSvc := service.NewUserService(userRepo, env.JwtSecret)
	userSvc = service.NewUserServiceSQS(userSvc, sqsAWSConfig, env.SQSUserLogsQueueURL)
	userController := transport.NewUserController(r, userSvc, env.JwtSecret)
	userController.RegisterRoutes()

	userLogRepo := repo.NewLogRepo(dynamodbAWSConfig, env.DynamoTable)
	adminSvc := service.NewAdminService(userRepo, userLogRepo)
	adminControler := transport.NewAdminController(r, userSvc, adminSvc, env.JwtSecret)
	adminControler.RegisterRoutes()

	port := fmt.Sprintf(":%d", env.Port)
	srv := &http.Server{Addr: port, Handler: r}
	log.Println("API listening on " + port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
