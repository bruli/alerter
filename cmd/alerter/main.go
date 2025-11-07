package main

import (
	"context"
	"errors"
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
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log := buildLogger()

	conf, err := config.New()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load configuration")
		os.Exit(1)
	}

	publisher, err := telegram.NewPublisher(conf.TelegramToken, conf.TelegramChatID)
	if err != nil {
		log.Err(err).Msg("failed to create telegram publisher")
		os.Exit(1)
	}
	log.Info().Msgf("telegram publisher created")

	cache := memory.NewCache(time.Minute)
	go cache.RunTTL(ctx, 20*time.Second)

	publisSvc := message.NewPublish(publisher, cache)

	consumer := infranats.NewMessageConsumer(publisSvc)

	natsUrl := conf.NatsServerURL
	nc, err := nats.Connect(natsUrl)
	if err != nil {
		log.Err(err).Msg("failed to connect to nats")
		os.Exit(1)
	}
	log.Info().Msgf("connected to NATS server %s", natsUrl)
	defer func() {
		_ = nc.Drain()
		nc.Close()
	}()

	subject := infranats.PingSubject

	sub, err := nc.Subscribe(subject, consumer.Consume)
	if err != nil {
		log.Err(err).Msg("failed to subscribe to subject")
	}
	defer func() {
		_ = sub.Unsubscribe()
	}()

	log.Info().Msgf("subscribed to subject %s", subject)

	srv := infrahttp.NewServer(log)

	go shutDown(ctx, log, srv)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("error while starting server")
	}
}

func buildLogger() *zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &log
}

func shutDown(ctx context.Context, log *zerolog.Logger, srv *http.Server) {
	<-ctx.Done()
	log.Info().Msg("shutdown signal received, stopping server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("error shutting down HTTP server")
	}
}
