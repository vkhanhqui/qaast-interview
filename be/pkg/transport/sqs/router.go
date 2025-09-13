package sqs

import (
	"be/pkg/log"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func NewSQSRouter(routeFn func(msg types.Message) string) Router {
	return &sqsRouter{
		routeFn:     routeFn,
		middlewares: []Middleware{},
		handlers:    map[string]HandlerFunc{},
	}
}

func RouteFromAttributeFn(msg types.Message) (route string) {
	attr, ok := msg.MessageAttributes["route"]
	if !ok {
		return
	}

	route = *attr.StringValue
	return
}

type sqsRouter struct {
	routeFn     func(msg types.Message) string // func used to extract the route name from the message
	middlewares []Middleware
	handlers    map[string]HandlerFunc
}

func (r *sqsRouter) AddHandler(route string, handler HandlerFunc) {
	r.handlers[route] = handler
}

func (r *sqsRouter) Use(mws ...Middleware) {
	if len(r.handlers) > 0 {
		panic("Messaging Router: all middlewares must be defined before adding any handler")
	}
	r.middlewares = append(r.middlewares, mws...)
}

func (r *sqsRouter) Handle(ctx context.Context, msg types.Message) error {
	route := r.routeFn(msg)
	hf, ok := r.handlers[route]
	if !ok {
		return r.NotFoundHandler(ctx, route)
	}

	// execute middlewares stack
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		currentMiddleware := r.middlewares[i]
		hf = currentMiddleware(hf)
	}

	return hf(ctx, msg)
}

func (r *sqsRouter) NotFoundHandler(ctx context.Context, route string) error {
	log.Info(fmt.Sprintf("Handler was not found for route %s", route))
	return nil
}
