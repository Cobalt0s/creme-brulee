package messaging

import (
	"github.com/google/uuid"
	"time"
)

type Baggage struct {
	ActorID     *uuid.UUID   `json:"actorId,omitempty"`
	ReceiverIDs []uuid.UUID  `json:"receiverIds,omitempty"`
	SentAt      TimeNano3339 `json:"sentAt,omitempty"`
}

func (b *Baggage) ResendTo(receivers []uuid.UUID) *Baggage {
	b.ReceiverIDs = receivers
	b.SentAt = NewNano3339Time(time.Now())
	return b
}

type BFactory struct {
	*Baggage
}

func BaggageFactory() BFactory {
	return BFactory{&Baggage{}}
}

func (f *BFactory) By(actorID uuid.UUID) *BFactory {
	f.Baggage.ActorID = &actorID
	return f
}

func (f BFactory) StampForOne(receiver uuid.UUID) Baggage {
	return f.StampFor([]uuid.UUID{receiver})
}

func (f BFactory) StampFor(receivers []uuid.UUID) Baggage {
	f.ResendTo(receivers)
	return *f.Baggage
}
