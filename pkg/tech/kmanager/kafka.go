package kmanager

import (
	"context"
	"fmt"
	config2 "github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

var (
	DefaultTopics = []string{DeadLetterQueueTopic}
)

func CreateTopics(ctx context.Context, baseConf *config2.BaseConfig, cfg *config2.KafkaConfig, topicNames []string) {
	log := ctxlogrus.Extract(ctx)

	if baseConf.Env != "testing" {
		return
	}

	ka, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": cfg.Host,
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to init kafka admin client %v", err))
	}

	topicNames = addDefaultTopics(topicNames)

	topics := make([]kafka.TopicSpecification, len(topicNames))
	for i, name := range topicNames {
		topics[i] = kafka.TopicSpecification{
			Topic:             name,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	_, err = ka.CreateTopics(ctx, topics)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to auto create topics %v", err))
	}
}

func addDefaultTopics(topicNames []string) []string {
	// use set for unique topics
	uniqueTopics := make(map[string]bool)
	for _, name := range topicNames {
		uniqueTopics[name] = true
	}
	for _, name := range DefaultTopics {
		uniqueTopics[name] = true
	}
	// map to slice conversion
	result := make([]string, len(uniqueTopics))
	index := 0
	for topic := range uniqueTopics {
		result[index] = topic
		index += 1
	}
	return topicNames
}
