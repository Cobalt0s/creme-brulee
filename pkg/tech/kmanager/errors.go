package kmanager

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func NewKafkaErrEventConsume(message *kafka.Message, topic string, err error) error {
	return fmt.Errorf("couldn't consume event from %v topic due to %v", topic, err)
}

func NewKafkaErrEventParse(message *kafka.Message, topic string, err error) error {
	return fmt.Errorf("couldn't parse event %v from %v topic due to %v", string(message.Value), topic, err)
}

func NewKafkaErrUnknownTopic(topic string) error {
	return fmt.Errorf("incoming message is from unknown topic %v", topic)
}
