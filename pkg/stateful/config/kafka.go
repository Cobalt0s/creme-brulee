package config

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConfig struct {
	Host     string
	Username string
	Password string
	SslMode  string
}

func NewKafkaConfig() (*KafkaConfig, error) {

	envKafkaHost, err := GetEnv("KAFKA_HOST")
	if err != nil {
		return nil, err
	}
	envKafkaSslMode, err := GetEnv("KAFKA_SSL")
	if err != nil {
		return nil, err
	}

	var envKafkaUsername string
	var envKafkaPassword string
	if envKafkaSslMode != "disable" {
		envKafkaUsername, err = GetEnv("KAFKA_USERNAME")
		if err != nil {
			return nil, err
		}
		envKafkaPassword, err = GetEnv("KAFKA_PASSWORD")
		if err != nil {
			return nil, err
		}
	}

	return &KafkaConfig{
		Host:     envKafkaHost,
		Username: envKafkaUsername,
		Password: envKafkaPassword,
		SslMode:  envKafkaSslMode,
	}, nil
}

func (cfg *KafkaConfig) GetKafkaConfigMapConsumer(clientID, consumerGroup string) *kafka.ConfigMap {
	result := *cfg.getKafkaConfigMapShared()
	result["client.id"] = clientID
	result["group.id"] = consumerGroup
	result["enable.auto.commit"] = false
	return &result
}

func (cfg *KafkaConfig) GetKafkaConfigMapConsumerEarliest(clientID, consumerGroup string) *kafka.ConfigMap {
	result := *cfg.GetKafkaConfigMapConsumer(consumerGroup, clientID)
	result["auto.offset.reset"] = "earliest"
	return &result
}

func (cfg *KafkaConfig) GetKafkaConfigMapProducer(clientID string) *kafka.ConfigMap {
	result := *cfg.getKafkaConfigMapShared()
	result["client.id"] = clientID
	return &result
}

func (cfg *KafkaConfig) getKafkaConfigMapShared() *kafka.ConfigMap {
	if cfg.SslMode != "disable" {
		return &kafka.ConfigMap{
			"bootstrap.servers": cfg.Host,
			"security.protocol": "SASL_SSL",
			"sasl.mechanisms":   "PLAIN",
			"sasl.username":     cfg.Username,
			"sasl.password":     cfg.Password,
		}
	}

	return &kafka.ConfigMap{
		"bootstrap.servers": cfg.Host,
	}
}
