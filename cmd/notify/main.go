package main

import (
	"context"
	"github.com/patyukin/mbs-notify/internal/config"
	"github.com/patyukin/mbs-notify/internal/telegram"
	"github.com/patyukin/mbs-notify/internal/usecase"
	"github.com/patyukin/mbs-pkg/pkg/rabbitmq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
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
		rabbitmq.AuthNotifyQueue,
		[]string{rabbitmq.AuthSignInConfirmCodeRouteKey, rabbitmq.AuthSignUpResultMessageRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind AuthNotifyQueue to exchange with - AuthSignUpResultMessageRouteKey, "+
			"AuthSignInConfirmCodeRouteKey: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.PaymentNotifyQueue,
		[]string{rabbitmq.AccountCreationRouteKey, rabbitmq.PaymentExecutionInitiateRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind PaymentNotifyQueue to exchange with - PaymentExecutionProcessRouteKey: %v", err)
	}

	err = rbt.BindQueueToExchange(
		rabbitmq.Exchange,
		rabbitmq.NotifyAuthQueue,
		[]string{rabbitmq.NotifySignUpConfirmCodeRouteKey},
	)
	if err != nil {
		log.Fatal().Msgf("failed to bind NotifyAuthQueue to exchange with - NotifySignUpConfirmCodeRouteKey: %v", err)
	}

	uc := usecase.New(tg, rbt)

	go uc.StartTelegramBot(ctx)
	go func() {
		if err = rbt.Consume(ctx, rabbitmq.AuthNotifyQueue, uc.AuthConsumeHandler); err != nil {
			log.Fatal().Msgf("failed to start auth_notify_consumer: %v", err)
		}
	}()

	go func() {
		if err = rbt.Consume(ctx, rabbitmq.PaymentNotifyQueue, uc.PaymentConsumeHandler); err != nil {
			log.Fatal().Msgf("failed to start auth_notify_consumer: %v", err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	log.Info().Msg("Завершение работы...")
}
