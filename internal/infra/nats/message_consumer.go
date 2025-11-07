package nats

import (
	"context"
	"encoding/json"

	"github.com/bruli/alerter/internal/domain/message"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

const PingSubject = "ping.created"

type Message struct {
	Resource string
	Status   string
}

type MessageConsumer struct {
	svc *message.Publish
	log *zerolog.Logger
}

func (c MessageConsumer) Consume(msg *nats.Msg) {
	var m Message
	if err := json.Unmarshal(msg.Data, &m); err != nil {
		c.log.Error().Err(err).Msg("error while unmarshalling message")
		return
	}
	ms, err := message.NewMessage(m.Resource, m.Status)
	if err != nil {
		c.log.Err(err).Msg("fail to create a message")
		return
	}
	if err = c.svc.Handle(context.Background(), ms); err != nil {
		c.log.Err(err).Msg("fail to handle message")
	}
}

func NewMessageConsumer(svc *message.Publish, log *zerolog.Logger) *MessageConsumer {
	return &MessageConsumer{svc: svc, log: log}
}
