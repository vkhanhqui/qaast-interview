package sqs

import (
	"be/pkg/log"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	defaultMaxMessages       = 1
	defaultVisibilityTimeout = 300
	defaultWaitTimeSeconds   = 10
)

type Config struct {
	AWSConfig aws.Config
	QueueURL  string
	// MaxMessages is also used as channel buffer capacity.
	// So it should be match with worker capacity ("drop" pattern) to avoid duplication
	// in case the message backs to queue after timed-out.
	MaxMessages       int
	VisibilityTimeout int
	WaitTimeSeconds   int

	NAckVisibilityTimeout int
}

func (c *Config) SetDefaultValues() {
	if c.MaxMessages == 0 {
		c.MaxMessages = defaultMaxMessages
	}
	if c.VisibilityTimeout == 0 {
		c.VisibilityTimeout = defaultVisibilityTimeout
	}
	if c.NAckVisibilityTimeout == 0 {
		c.NAckVisibilityTimeout = defaultVisibilityTimeout
	}
	if c.WaitTimeSeconds == 0 {
		c.WaitTimeSeconds = defaultWaitTimeSeconds
	}
}

type Worker interface {
	Run(context.Context) error
	Stop(context.Context) error
}

type Router interface {
	AddHandler(route string, handler HandlerFunc)
	Use(mws ...Middleware)
	Handle(ctx context.Context, msg types.Message) error
}

type HandlerFunc func(ctx context.Context, msg types.Message) error
type Middleware func(next HandlerFunc) HandlerFunc

type BatchRouter interface {
	AddHandler(route string, handler BatchHandlerFunc)
	Use(mws ...BatchMiddleware)
	Handle(ctx context.Context, msgs []types.Message) error
}

type BatchHandlerFunc func(ctx context.Context, msgs []types.Message) error
type BatchMiddleware func(next BatchHandlerFunc) BatchHandlerFunc

func ListenForTermination(ctx context.Context, w Worker, l log.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-stop
		l.Info("Received termination signal, stopping the worker")

		err := w.Stop(ctx)
		if err != nil {
			l.Error(err.Error(), err)
		}
	}()
}
