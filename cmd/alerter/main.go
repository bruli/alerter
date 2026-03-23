package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bruli/alerter/internal/config"
	"github.com/bruli/alerter/internal/domain/message"
	infrahttp "github.com/bruli/alerter/internal/infra/http"
	"github.com/bruli/alerter/internal/infra/memory"
	infranats "github.com/bruli/alerter/internal/infra/nats"
	"github.com/bruli/alerter/internal/infra/telegram"
	"github.com/bruli/pinger/pkg/events"
	"github.com/nats-io/nats.go"
)

const serviceName = "alerter"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := buildLog()

	conf, err := config.New()
	if err != nil {
		log.ErrorContext(ctx, "failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	publisher, err := telegram.NewPublisher(conf.TelegramToken, conf.TelegramChatID)
	if err != nil {
		log.ErrorContext(ctx, "failed to create telegram publisher", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.InfoContext(ctx, "telegram publisher created")

	cache := memory.NewCache(conf.CacheTTL)
	go cache.RunTTL(ctx, conf.CacheTTLRefresh)

	publisSvc := message.NewPublish(publisher, cache)

	consumer := infranats.NewPingResultConsumer(publisSvc, log)

	natsUrl := conf.NatsServerURL
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.ErrorContext(ctx, "failed to connect to nats", slog.String("error", err.Error()))
		os.Exit(1)
	}
	log.InfoContext(ctx, "nats connected", slog.String("url", natsUrl))
	defer func() {
		_ = nc.Drain()
		nc.Close()
	}()

	subject := events.PingSubjet

	sub, err := nc.Subscribe(subject, consumer.Consume)
	if err != nil {
		log.ErrorContext(ctx, "failed to subscribe to subject",
			slog.String("subject", subject),
			slog.String("error", err.Error()))
	}
	defer func() {
		_ = sub.Unsubscribe()
	}()
	log.InfoContext(ctx, "nats subscription created", slog.String("subject", subject))

	srv := infrahttp.NewServer(log)

	go shutDown(ctx, log, srv)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.ErrorContext(ctx, "failed to start server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func buildLog() *slog.Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	log := slog.New(handler)
	log.With("service", serviceName)
	return log
}

func shutDown(ctx context.Context, log *slog.Logger, srv *http.Server) {
	<-ctx.Done()
	log.InfoContext(ctx, "shutdown signal received, stopping server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.ErrorContext(ctx, "error shutting down HTTP server", slog.String("error", err.Error()))
	}
}
