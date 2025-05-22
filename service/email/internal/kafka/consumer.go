package kafka

import "github.com/IBM/sarama"

func NewConsumer(brokers []string, groupID string, topics []string, handler sarama.ConsumerGroupHandler) error {
	config := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return err
	}

	go func() {
		for {
			err := consumerGroup.Consume(nil, topics, handler)
			if err != nil {
				panic(err)
			}
		}
	}()

	return nil
}
