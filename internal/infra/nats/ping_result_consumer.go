package nats

import (
	"context"

	"github.com/bruli/alerter/internal/domain/message"
	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

type PingResultConsumer struct {
	svc *message.Publish
	log *zerolog.Logger
}

func (c PingResultConsumer) Consume(msg *nats.Msg) {
	var m events.PingResult
	if err := proto.Unmarshal(msg.Data, &m); err != nil {
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

func NewPingResultConsumer(svc *message.Publish, log *zerolog.Logger) *PingResultConsumer {
	return &PingResultConsumer{svc: svc, log: log}
}
