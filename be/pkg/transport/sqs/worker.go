package sqs

import (
	"be/pkg/log"
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func NewWorker(c Config, r Router, l log.Logger) Worker {
	w := &worker{
		config: c,
		router: r,
		client: sqs.NewFromConfig(c.AWSConfig),
		logger: l,
	}

	w.config.SetDefaultValues()
	return w
}

type worker struct {
	client *sqs.Client
	config Config
	router Router
	stop   bool
	logger log.Logger
}

func (w *worker) Run(ctx context.Context) error {
	w.logger.Info("SQS worker is now running...")

	rmi := &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(w.config.QueueURL),
		MaxNumberOfMessages:   int32(w.config.MaxMessages),
		VisibilityTimeout:     int32(w.config.VisibilityTimeout),
		WaitTimeSeconds:       int32(w.config.WaitTimeSeconds),
		AttributeNames:        []types.QueueAttributeName{"ApproximateReceiveCount"},
		MessageAttributeNames: []string{"All"},
	}

	for {
		out, err := w.client.ReceiveMessage(ctx, rmi)
		if err != nil {
			return err
		}

		if w.stop {
			return nil
		}

		var wg sync.WaitGroup
		for _, msg := range out.Messages {
			wg.Add(1)
			go w.processMessage(ctx, msg, &wg)
		}
		wg.Wait()
	}
}

func (w *worker) Stop(ctx context.Context) (err error) {
	w.stop = true
	return
}

func (w *worker) processMessage(ctx context.Context, msg types.Message, wg *sync.WaitGroup) {
	defer log.PanicRecover(w.logger)

	jobID := *msg.MessageId
	idAttr, ok := msg.MessageAttributes["id"]
	if ok {
		jobID = aws.ToString(idAttr.StringValue)
	}

	w.logger.Info(fmt.Sprintf("Processing job %s", jobID))
	defer wg.Done()

	err := w.router.Handle(ctx, msg)
	success, rc := w.handleError(ctx, err, msg)

	if !success {
		if rc%3 == 0 {
			w.logger.Error(err.Error(), err, msg)
		} else {
			w.logger.Info(err.Error(), err, msg)
		}
		return
	}

	_, err = w.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(w.config.QueueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		w.logger.Error(err.Error(), err)
		return
	}

	w.logger.Info(fmt.Sprintf("Finished processing job %s", jobID))
}

func (w *worker) handleError(ctx context.Context, cause error, msg types.Message) (success bool, retryCount int) {
	if cause == nil {
		return true, 0
	}

	// handle retry
	rc, err := strconv.Atoi(msg.Attributes["ApproximateReceiveCount"])
	if err != nil {
		w.logger.Info("Could not parse message receive count: "+err.Error(), err, msg)
		return false, 0
	}

	_, err = w.client.ChangeMessageVisibility(ctx, &sqs.ChangeMessageVisibilityInput{
		QueueUrl:          aws.String(w.config.QueueURL),
		ReceiptHandle:     msg.ReceiptHandle,
		VisibilityTimeout: int32(max(w.config.VisibilityTimeout, 60) * rc),
	})
	if err != nil {
		w.logger.Info("Could not change message visibility timeout: "+err.Error(), err, msg)
	}

	return false, rc
}
