package kmanager

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/messaging"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"gorm.io/gorm"
	"time"
)

var (
	DeadLetterQueueTopic = "DLQ"
	retryTimeouts        = []int{0, 2, 4, 8, 16}
)

type MessageConsumer struct {
	consumer    *kafka.Consumer
	topicNames  []string
	db          *gorm.DB
	arguments   *messaging.ContextualArguments
	stop        chan bool
	dlqProducer *MessageProducer
}

type KafkaHealthChecker struct {
	consumer *kafka.Consumer
}

func NewMessageConsumer(ctx context.Context, arguments *messaging.ContextualArguments, db *gorm.DB, cfg *config.KafkaConfig, clientID, consumerGroup string, topicNames []string) *MessageConsumer {
	return createFailsafeMessageConsumer(ctx, arguments, db, cfg, clientID, consumerGroup, topicNames, false)
}

func NewMessageConsumerOrigin(ctx context.Context, arguments *messaging.ContextualArguments, db *gorm.DB, cfg *config.KafkaConfig, clientID, consumerGroup string, topicNames []string) *MessageConsumer {
	return createFailsafeMessageConsumer(ctx, arguments, db, cfg, clientID, consumerGroup, topicNames, true)
}

func (mc *MessageConsumer) Close() error {
	mc.stop <- true
	return mc.consumer.Close()
}

func createFailsafeMessageConsumer(ctx context.Context, arguments *messaging.ContextualArguments, db *gorm.DB, cfg *config.KafkaConfig, clientID, consumerGroup string, topicNames []string, origin bool) *MessageConsumer {
	log := ctxlogrus.Extract(ctx)

	kc, err := kafka.NewConsumer(getConsumerMap(origin, cfg, clientID, consumerGroup))
	if err != nil {
		log.Fatal(err)
	}
	dlqProducer, err := NewMessageProducer(cfg, clientID)
	if err != nil {
		log.Fatal(err)
	}

	return &MessageConsumer{
		consumer:    kc,
		topicNames:  topicNames,
		db:          db,
		arguments:   arguments,
		stop:        make(chan bool),
		dlqProducer: dlqProducer,
	}
}

func getConsumerMap(origin bool, cfg *config.KafkaConfig, clientID, consumerGroup string) *kafka.ConfigMap {
	if origin {
		return cfg.GetKafkaConfigMapConsumerEarliest(clientID, consumerGroup)
	}
	return cfg.GetKafkaConfigMapConsumer(clientID, consumerGroup)
}

type TopicHandler func(ctx context.Context, arguments *messaging.ContextualArguments,
	db *gorm.DB, msg *kafka.Message) error

func (mc *MessageConsumer) Start(ctx context.Context, handleMessage TopicHandler) error {
	log := ctxlogrus.Extract(ctx)
	log.Info("starting kafka consumer")

	if err := mc.consumer.SubscribeTopics(mc.topicNames, nil); err != nil {
		log.Errorf("failed to subscirbe to kafka topics %v", err)
		return err
	}

	for {
		select {
		case <-mc.stop:
			log.Info("stopping message consumer")
			return nil
		default:
			if msg, err := mc.consumer.ReadMessage(5 * time.Second); err == nil {
				if len(msg.Key) != 0 { // ignore heartbeat messages
					for _, retryTimeout := range retryTimeouts {
						waitTime := time.Duration(retryTimeout) * time.Second
						time.Sleep(waitTime)
						err := handleMessage(ctx, mc.arguments, mc.db, msg)
						if err != nil {
							if retryTimeout != retryTimeouts[len(retryTimeouts)-1] {
								log.Error("failed processing event, will retry")
							} else {
								log.Error("failed processing event, reached final retry")
								if dlqErr := mc.dlqProducer.ProduceMessage(ctx, &DLQMessage{
									MessageValue: msg.Value,
									FormerTopic:  msg.TopicPartition.Topic,
								}); dlqErr != nil {
									return dlqErr
								}
							}
						} else {
							break
						}
					}
					if _, err := mc.consumer.CommitMessage(msg); err != nil {
						log.Warnf("couldn't commit message %v", err)
					}
				}
			} else {
				if kErr, ok := err.(kafka.Error); ok {
					if kErr.Code() == kafka.ErrTimedOut {
						continue
					}
				}
				log.Warnf("consumer kafka error: %v (%v)\n", err, msg)
			}

		}
	}
}
