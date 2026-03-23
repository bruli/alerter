package nats

import (
	"context"
	"log/slog"

	"github.com/bruli/alerter/internal/domain/message"
	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type PingResultConsumer struct {
	svc *message.Publish
	log *slog.Logger
}

func (c PingResultConsumer) Consume(msg *nats.Msg) {
	var m events.PingResult
	ctx := context.Background()
	if err := proto.Unmarshal(msg.Data, &m); err != nil {
		c.log.ErrorContext(ctx, "failed to unmarshal message")
		return
	}
	ms, err := message.NewMessage(m.Resource, m.Status)
	if err != nil {
		c.log.ErrorContext(ctx, "fail to create a message", slog.String("error", err.Error()))
		return
	}
	if err = c.svc.Handle(context.Background(), ms); err != nil {
		c.log.ErrorContext(ctx, "fail to handle message", slog.String("error", err.Error()))
	}
}

func NewPingResultConsumer(svc *message.Publish, log *slog.Logger) *PingResultConsumer {
	return &PingResultConsumer{svc: svc, log: log}
}
