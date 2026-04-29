package kafka

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"go.uber.org/zap"
)

type Kafka struct {
	logger   logger.LoggerInterface
	producer sarama.SyncProducer
	brokers  []string
}

func NewKafka(logger logger.LoggerInterface, brokers []string) *Kafka {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		logger.Fatal("Failed to create Kafka producer", zap.Error(err))
	}

	logger.Info("Kafka producer connected successfully")

	return &Kafka{
		producer: producer,
		brokers:  brokers,
		logger:   logger,
	}
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
		return err
	}

	k.logger.Info("Message is stored in topic", zap.String("topic", topic), zap.Int32("partition", partition), zap.Int64("offset", offset))

	return nil
}

func (k *Kafka) StartConsumers(topics []string, groupID string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumerGroup, err := sarama.NewConsumerGroup(k.brokers, groupID, config)
	if err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, handler); err != nil {
				k.logger.Error("Error from consumer", zap.Error(err))
			}
		}
	}()

	go func() {
		for err := range consumerGroup.Errors() {
			k.logger.Error("Consumer group error", zap.Error(err))
		}
	}()

	return nil
}
