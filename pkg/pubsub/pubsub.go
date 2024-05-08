package pubsub

import (
	"bulk/config"
	"bulk/logger"
	"bulk/tracer"
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"cloud.google.com/go/pubsub"

	"go.opentelemetry.io/otel/attribute"
)

var wg sync.WaitGroup

type PubSub struct {
	l      logger.ILogger
	client *pubsub.Client
}

//go:generate mockgen -source pubsub.go -destination pubsub_mock.go -package pubsub

type PubSubClient interface {
	PublishData(ctx context.Context, topicName string, data []byte, attributes map[string]string)
}

func NewPubSubClient(logger logger.ILogger, ctx context.Context, conf *config.Config) *PubSub {
	childCtx, span := tracer.StartSpan(ctx, "NewPubSubClient")
	defer span.End()

	pubsubClient, err := pubsub.NewClient(childCtx, conf.ProjectId)
	if err != nil {
		span.RecordError(err)
		logger.Errorf("Error while creating new pubsub client %v", err)
		return nil
	}
	return &PubSub{l: logger, client: pubsubClient}
}

func (p *PubSub) PublishData(ctx context.Context, topicName string, data []byte, attributes map[string]string) {
	childCtx, span := tracer.StartSpan(ctx, "PublishData")
	defer span.End()
	span.SetAttributes(
		attribute.String("topicName", topicName),
		attribute.String("data", string(data)),
	)

	t := p.client.Topic(topicName)
	p.l.Infof("published data to %v", topicName)
	result := t.Publish(childCtx, &pubsub.Message{
		Data:       data,
		Attributes: attributes,
	})
	var totalErrors uint64
	wg.Add(1)
	go func(result *pubsub.PublishResult) {
		defer wg.Done()
		// The Get method blocks until a server-generated ID or
		// an error is returned for the published message.
		_, err := result.Get(childCtx)
		if err != nil {
			// Error handling code can be added here.
			p.l.Errorf("Failed to publish %v", err)
			atomic.AddUint64(&totalErrors, 1)
			return
		}
	}(result)
	wg.Wait()

	if totalErrors > 0 {
		p.l.Error(fmt.Sprintf("%d messages did not publish successfully", totalErrors))
		return
	}
}
