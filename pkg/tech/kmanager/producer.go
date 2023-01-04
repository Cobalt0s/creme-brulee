package kmanager

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/messaging"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/logging"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"go.opentelemetry.io/otel/codes"
	"gorm.io/gorm"
)

func SendEvent(ctx context.Context, tx *gorm.DB, topic string, event messaging.JSONConvertable) error {
	ctx, span := logging.StartSpan(ctx, "SendEvent")
	defer span.End()
	log := ctxlogrus.Extract(ctx)

	jsonData, err := event.ToJSON()
	if err != nil {
		return err
	}
	if err = tx.Create(&OutboxORM{
		KafkaTopic: topic,
		KafkaKey:   "some-kafka-key",
		KafkaValue: jsonData,
	}).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func SendEvents(ctx context.Context, tx *gorm.DB, topic string, events []messaging.JSONConvertable) error {
	if len(events) == 0 {
		// no events to send, no op
		return nil
	}

	ctx, span := logging.StartSpan(ctx, "SendEvents")
	defer span.End()
	log := ctxlogrus.Extract(ctx)

	outboxORMs := make([]OutboxORM, len(events))
	for i, e := range events {
		jsonData, err := e.ToJSON()
		if err != nil {
			return err
		}
		outboxORMs[i] = OutboxORM{
			KafkaTopic: topic,
			KafkaKey:   "some-kafka-key",
			KafkaValue: jsonData,
		}
	}

	if err := tx.Create(outboxORMs).Error; err != nil {
		log.Error(err)
		return err
	}
	return nil
}

type MessageProducer struct {
	producer *kafka.Producer
}

func NewMessageProducer(cfg *config.KafkaConfig, clientID string) (*MessageProducer, error) {
	kafkaConfigMap := cfg.GetKafkaConfigMapProducer(clientID)
	producer, err := kafka.NewProducer(kafkaConfigMap)
	if err != nil {
		return nil, err
	}
	return &MessageProducer{producer: producer}, nil
}

func (p *MessageProducer) ProduceMessage(ctx context.Context, message Message) error {
	value := message.Key()
	key := message.Value()

	ctx, span := TraceFromEventNamed(ctx, value, "Harvested")
	defer span.End()
	log := ctxlogrus.Extract(ctx)

	delivery := make(chan kafka.Event, 1)

	err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: message.Topic(), Partition: kafka.PartitionAny},
		Key:            key,
		Value:          enhanceWithCurrentTrace(ctx, value),
	},
		delivery,
	)
	if err != nil {
		log.Error("failed to enqueue message for kafka producer")
		return err
	}

	output := <-delivery
	messageReport := output.(*kafka.Message)

	err = messageReport.TopicPartition.Error
	if err != nil {
		log.Error("failed sending message")
		span.SetStatus(codes.Error, "failed sending message")
		return err
	} else {
		log.Debugf("Message delivered %v", string(messageReport.Value))
	}

	return nil
}
