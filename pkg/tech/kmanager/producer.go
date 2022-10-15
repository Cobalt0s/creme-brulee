package kmanager

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/messaging"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"gorm.io/gorm"
)

type OutboxORM struct {
	ID         uint   `gorm:"primarykey"`
	KafkaTopic string `gorm:"type:text"`
	KafkaKey   string `gorm:"type:text"`
	KafkaValue string `gorm:"type:text"`
}

func (*OutboxORM) TableName() string {
	return "outbox"
}

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
