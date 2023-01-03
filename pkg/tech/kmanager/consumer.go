package kmanager

import (
	"context"
	"github.com/Cobalt0s/creme-brulee/pkg/rest/messaging"
	"github.com/Cobalt0s/creme-brulee/pkg/stateful/config"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"gorm.io/gorm"
)

type MessageConsumer struct {
	consumer   *kafka.Consumer
	topicNames []string
	db         *gorm.DB
	arguments  *messaging.ContextualArguments
	done       chan bool
	poller     chan *kafka.Message
}

type KafkaHealthChecker struct {
	consumer *kafka.Consumer
}

func NewMessageConsumer(ctx context.Context, arguments *messaging.ContextualArguments, db *gorm.DB, cfg *config.KafkaConfig, clientID, consumerGroup string, topicNames []string) *MessageConsumer {
	return createMessageConsumer(ctx, arguments, db, cfg, clientID, consumerGroup, topicNames, false)
}

func NewMessageConsumerOrigin(ctx context.Context, arguments *messaging.ContextualArguments, db *gorm.DB, cfg *config.KafkaConfig, clientID, consumerGroup string, topicNames []string) *MessageConsumer {
	return createMessageConsumer(ctx, arguments, db, cfg, clientID, consumerGroup, topicNames, true)
}

func (mc *MessageConsumer) Close() error {
	select {
	case mc.done <- true:
	default:
	}
	return mc.consumer.Close()
}

func createMessageConsumer(ctx context.Context, arguments *messaging.ContextualArguments, db *gorm.DB, cfg *config.KafkaConfig, clientID, consumerGroup string, topicNames []string, origin bool) *MessageConsumer {
	log := ctxlogrus.Extract(ctx)

	kc, kafkaError := kafka.NewConsumer(getConsumerMap(origin, cfg, clientID, consumerGroup))
	if kafkaError != nil {
		log.Fatal(kafkaError)
	}
	return &MessageConsumer{
		arguments:  arguments,
		consumer:   kc,
		topicNames: topicNames,
		db:         db,
		done:       make(chan bool),
		poller:     make(chan *kafka.Message),
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

	go func() {
		// infinite Poller
		if msg, err := mc.consumer.ReadMessage(-1); err == nil {
			if len(msg.Key) != 0 { // ignore heartbeat messages
				mc.poller <-msg
			}
		} else {
			log.Warnf("consumer kafka error: %v (%v)\n", err, msg)
		}
	}()

	for {
		select {
		case <-mc.done:
			return nil
		case msg := <-mc.poller:
			if err := handleMessage(ctx, mc.arguments, mc.db, msg); err != nil {
				log.Error("message will be NOT committed")
				// TODO if we couldn't consume event from kafka we need to retry
				// TODO do we fail the pod since it couldn't consume event?
				continue
			}
			if _, err := mc.consumer.CommitMessage(msg); err != nil {
				log.Warnf("couldn't commit message %v", err)
			}
		}
	}
}
