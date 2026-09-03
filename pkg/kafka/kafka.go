package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
)

type Kafka struct {
	logger   logger.LoggerInterface
	producer sarama.SyncProducer
	brokers  []string

	publishedMessages metric.Int64Counter
	publishErrors     metric.Int64Counter
	consumeErrors     metric.Int64Counter
}

func newKafkaMetrics() (*Kafka, error) {
	meter := otel.Meter("kafka")

	publishedMessages, err := meter.Int64Counter(
		"kafka_published_messages_total",
		metric.WithDescription("Total number of messages successfully published to Kafka"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	publishErrors, err := meter.Int64Counter(
		"kafka_publish_errors_total",
		metric.WithDescription("Total number of Kafka publish failures"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	consumeErrors, err := meter.Int64Counter(
		"kafka_consume_errors_total",
		metric.WithDescription("Total number of Kafka consumer errors"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		publishedMessages: publishedMessages,
		publishErrors:     publishErrors,
		consumeErrors:     consumeErrors,
	}, nil
}

func NewKafka(logger logger.LoggerInterface, brokers []string) *Kafka {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	k := &Kafka{
		brokers: brokers,
		logger:  logger,
	}

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Kafka is an optional integration for local E2E runs. Keep the service
		// available when the broker is intentionally not started; SendMessage
		// and consumer startup below become safe no-ops in this mode.
		logger.Warn("Kafka unavailable; continuing with Kafka disabled", zap.Error(err))
		return k
	}

	k.producer = producer
	logger.Info("Kafka producer connected successfully")

	metricsK, err := newKafkaMetrics()
	if err != nil {
		logger.Warn("Failed to initialize Kafka metrics", zap.Error(err))
	} else if metricsK != nil {
		k.publishedMessages = metricsK.publishedMessages
		k.publishErrors = metricsK.publishErrors
		k.consumeErrors = metricsK.consumeErrors
	}

	return k
}

func (k *Kafka) SendMessage(topic string, key string, value []byte) error {
	if k == nil || k.producer == nil {
		return nil
	}
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		if k.publishErrors != nil {
			k.publishErrors.Add(context.Background(), 1, metric.WithAttributes(attribute.String("topic", topic)))
		}
		return err
	}

	if k.publishedMessages != nil {
		k.publishedMessages.Add(context.Background(), 1, metric.WithAttributes(attribute.String("topic", topic)))
	}

	k.logger.Info("Message is stored in topic", zap.String("topic", topic), zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}

// SendMessageWithHeaders publishes a message with explicit record headers in
// addition to the payload, propagating the active trace context as
// traceparent/tracestate headers. Used by the email service to attach retry/DLQ
// metadata (Phase 4) while keeping the payload an unchanged envelope.
func (k *Kafka) SendMessageWithHeaders(ctx context.Context, topic string, key string, value []byte, headers []sarama.RecordHeader) error {
	if k == nil || k.producer == nil {
		return nil
	}
	msg := &sarama.ProducerMessage{
		Topic:   topic,
		Key:     sarama.StringEncoder(key),
		Value:   sarama.ByteEncoder(value),
		Headers: headers,
	}
	if ctx != nil {
		carrier := propagation.MapCarrier{}
		otel.GetTextMapPropagator().Inject(ctx, carrier)
		if len(carrier) > 0 {
			for kk, vv := range carrier {
				msg.Headers = append(msg.Headers, sarama.RecordHeader{Key: []byte(kk), Value: []byte(vv)})
			}
		}
	}

	partition, offset, err := k.producer.SendMessage(msg)
	if err != nil {
		if k.publishErrors != nil {
			k.publishErrors.Add(context.Background(), 1, metric.WithAttributes(attribute.String("topic", topic)))
		}
		return err
	}

	if k.publishedMessages != nil {
		k.publishedMessages.Add(context.Background(), 1, metric.WithAttributes(attribute.String("topic", topic)))
	}

	k.logger.Info("Message is stored in topic", zap.String("topic", topic), zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}

// Close releases the Kafka producer. Consumers started with
// StartConsumersWithContext* stop when their context is cancelled (the consume
// goroutine defers consumerGroup.Close()); Close is idempotent and safe to call
// more than once.
func (k *Kafka) Close() error {
	if k == nil || k.producer == nil {
		return nil
	}
	err := k.producer.Close()
	k.producer = nil
	return err
}

func (k *Kafka) StartConsumers(topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	return k.StartConsumersWithContext(context.Background(), topics, groupID, handler)
}

// StartConsumersWithContext starts a consumer group with caller-owned
// cancellation. OffsetEarliest makes a new group replay retained events rather
// than silently starting at the newest offset.
func (k *Kafka) StartConsumersWithContext(ctx context.Context, topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	return k.startConsumers(ctx, topics, groupID, handler, true)
}

// StartConsumersWithContextManualCommit starts a consumer group with Sarama
// auto-commit disabled. The handler must call session.Commit() explicitly
// after reaching a terminal outcome (e.g. an email was sent or the message was
// rejected as invalid/duplicate). Messages that fail transiently are left
// uncommitted and are redelivered on the next rebalance (at-least-once).
func (k *Kafka) StartConsumersWithContextManualCommit(ctx context.Context, topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	return k.startConsumers(ctx, topics, groupID, handler, false)
}

func (k *Kafka) startConsumers(ctx context.Context, topics []string, groupID string, handler sarama.ConsumerGroupHandler, autoCommit bool) error {
	if k == nil || k.producer == nil {
		return nil
	}
	if ctx == nil {
		ctx = context.Background()
	}
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Offsets.AutoCommit.Enable = autoCommit

	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, groupID, config)
	if err != nil {
		return err
	}

	go func() {
		defer consumerGroup.Close()
		for {
			if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
				if ctx.Err() != nil {
					return
				}
				k.logger.Error("Error from consumer", zap.Error(err))
				if k.consumeErrors != nil {
					k.consumeErrors.Add(context.Background(), 1, metric.WithAttributes(attribute.String("group", groupID)))
				}
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case err, ok := <-consumerGroup.Errors():
				if !ok {
					return
				}
				k.logger.Error("Consumer group error", zap.Error(err))
				if k.consumeErrors != nil {
					k.consumeErrors.Add(context.Background(), 1, metric.WithAttributes(attribute.String("group", groupID)))
				}
			}
		}
	}()

	return nil
}
