package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/patyukin/mbs-notify/internal/config"
	"github.com/patyukin/mbs-notify/internal/telegram"
	"github.com/patyukin/mbs-notify/internal/usecase"
	"github.com/patyukin/mbs-pkg/pkg/kafka"
	"github.com/patyukin/mbs-pkg/pkg/mux_server"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal().Msgf("failed to load config, error: %v", err)
	}

	tg, err := telegram.New(cfg.TelegramToken)
	if err != nil {
		log.Fatal().Msgf("failed to create telegram: %v", err)
	}

	rbt, err := rabbitmq.New(cfg.RabbitMQUrl, rabbitmq.Exchange)
	if err != nil {
		log.Fatal().Msgf("failed to create rabbitmq: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.TelegramMessageQueue,
		[]string{rabbitmq.TelegramMessageRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind TelegramMessageQueue to exchange with - TelegramMessageRouteKey: %v", err)
	}

	kfk, err := kafka.NewProducer(cfg.Kafka.Brokers)
	if err != nil {
		log.Fatal().Msgf("failed to create kafka consumer, err: %v", err)
	}

	uc := usecase.New(tg, kfk)

	// mux server
	m := mux_server.New()

	errCh := make(chan error)

	go uc.StartTelegramBot(ctx)
	go func() {
		if err = rbt.Consume(ctx, rabbitmq.TelegramMessageQueue, uc.ConsumeTelegramMessageQueue); err != nil {
			log.Fatal().Msgf("failed to start auth_notify_consumer: %v", err)
		}
	}()

	// metrics + pprof server
	go func() {
		if err = m.Run(cfg.HttpServer.Port); err != nil {
			log.Error().Msgf("Failed to serve Prometheus metrics: %v", err)
			errCh <- err
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err = <-errCh:
		log.Error().Msgf("Failed to run, err: %v", err)
	case res := <-sigChan:
		if res == syscall.SIGINT || res == syscall.SIGTERM {
			log.Info().Msg("Signal received")
		} else if res == syscall.SIGHUP {
			log.Info().Msg("Signal received")
		}
	}

	log.Info().Msg("Shutting Down")

	// stop pprof server
	if err = m.Shutdown(ctx); err != nil {
		log.Error().Msgf("failed to shutdown pprof server: %s", err.Error())
	}
}
