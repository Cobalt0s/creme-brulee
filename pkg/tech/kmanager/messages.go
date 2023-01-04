package kmanager


type Message interface {
	Key() []byte
	Value() []byte
	Topic() *string
}

type OutboxORM struct {
	ID         uint   `gorm:"primarykey"`
	KafkaTopic string `gorm:"type:text"`
	KafkaKey   string `gorm:"type:text"`
	KafkaValue string `gorm:"type:text"`
}

func (*OutboxORM) TableName() string {
	return "outbox"
}

func (o *OutboxORM) Key() []byte {
	return []byte(o.KafkaKey)
}
func (o *OutboxORM) Value() []byte {
	return []byte(o.KafkaValue)
}
func (o *OutboxORM) Topic() *string {
	return &o.KafkaTopic
}

type DLQMessage struct {
	MessageValue []byte
	FormerTopic *string
}

func (m *DLQMessage) Key() []byte {
	if m.FormerTopic == nil {
		return []byte(("unknown_topic"))
	}
	return []byte(*m.FormerTopic)
}

func (m *DLQMessage) Value() []byte {
	return m.MessageValue
}

func (m *DLQMessage) Topic() *string {
	return &DeadLetterQueueTopic
}
